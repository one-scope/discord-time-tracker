package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/one-scope/discord-time-tracker/internal/api"
	"github.com/one-scope/discord-time-tracker/internal/config"
	"github.com/one-scope/discord-time-tracker/internal/db"
	"github.com/robfig/cron/v3"
)

const PREFIX_COMMAND = "tracker"

func New(aConfig *config.DiscordBotConfig, aDB *db.PostgresDB) (*Bot, error) {
	// DiscordSession初期化
	tSession, tError := discordgo.New("Bot " + aConfig.DiscordBotToken)
	if tError != nil {
		return nil, tError
	}
	tSession.Identify.Intents = discordgo.IntentGuildMembers | discordgo.IntentGuildPresences | discordgo.IntentGuildVoiceStates | discordgo.IntentGuilds | discordgo.IntentGuildMessages

	// Cron初期化
	tCron := cron.New()
	// DataManager初期化
	tManager := &dataManager{
		UsersByID:                map[string]*db.User{},
		StatusesByID:             map[string][]*db.Statuslog{},
		DB:                       aDB,
		PreViusStatusLogByUserID: map[string]*db.Statuslog{},
	}
	// Bot初期化
	tBot := &Bot{
		Session:         tSession,
		Cron:            tCron,
		FlushTimingCron: aConfig.FlushTimingCron,
		ErrorChannel:    aConfig.DiscordErrorChannelID,
		DataManager:     tManager,
	}
	return tBot, nil
}
func (aBot *Bot) Start() error {
	// イベントハンドラ設定
	aBot.setEventHandlers()
	if tError := aBot.Session.Open(); tError != nil {
		return tError
	}

	// 定期的にデータをDBに書き込む
	if _, tError := aBot.Cron.AddFunc(aBot.FlushTimingCron, func() {
		if tError := aBot.DataManager.flushData(); tError != nil {
			aBot.onEventError(aBot.Session, fmt.Sprintf("failed to flush data: %s", tError))
		}
	}); tError != nil {
		return tError
	}
	aBot.Cron.Start()

	return nil
}
func (aBot *Bot) Close() error {
	aBot.Cron.Stop()
	return aBot.Session.Close()
}

// イベント発生時のエラー処理
func (aBot *Bot) onEventError(aSession *discordgo.Session, aErrorMessage string) {
	if aBot.ErrorChannel == "" {
		log.Println("Warning: channel id is empty")
		log.Println("Error: ", aErrorMessage)
		return
	}

	log.Println("Error: ", aErrorMessage)
	if _, tError := aSession.ChannelMessageSend(aBot.ErrorChannel, "Error: "+aErrorMessage); tError != nil {
		log.Println("Error: failed to send message to error channel")
		log.Println(tError)
		return
	}
}

func (aBot *Bot) setEventHandlers() {
	// ボットが起動したとき。
	aBot.onEvent(aBot.event)
	// 誰かがサーバーに参加したとき。
	aBot.onGuildMemberAdd(aBot.guildMemberAdd)
	//誰かのロールが変わったとき。
	aBot.onGuildMemberUpdate(aBot.guildMemberUpdate)
	// 誰かがサーバーから退出したとき。
	aBot.onGuildMemberRemove(aBot.guildMemberRemove)
	//誰かのオンラインになったとき
	aBot.onPresenceUpdate(aBot.presenceUpdate)
	// 誰かの音声通話が更新されたとき。// 接続、切断もこれ。切断時は ChannelID が空文字。
	aBot.onVoiceStateUpdate(aBot.voiceStateUpdate)
	// 誰かがメッセージを送信したとき。
	aBot.onMessageCreate(aBot.messageCreate)
	//ロールが作成されたとき。
	aBot.onGuildRoleCreate(aBot.roleCreate)
	//ロールが更新されたとき。
	aBot.onGuildRoleUpdate(aBot.roleUpdate)
	//ロールが削除されたとき。
	aBot.onGuildRoleDelete(aBot.roleDelete)
}

// 　何かのイベントが発生したとき
func (aBot *Bot) event(aSession *discordgo.Session, aEvent *discordgo.Event) {
	log.Println("Debug: Event: ", aEvent.Type)
	if aEvent.Type == "GUILD_CREATE" {
		aBot.guildCreate(aSession, aEvent)
	}
}

// Botが起動したとき
func (aBot *Bot) guildCreate(aSession *discordgo.Session, aEvent *discordgo.Event) {
	var tMapRawData map[string]interface{}
	json.Unmarshal(aEvent.RawData, &tMapRawData)
	//GuildIDを取得
	tGuildID, tOk := tMapRawData["id"].(string)
	if !tOk {
		aBot.onEventError(aSession, "guildCreate: failed to convert guild id")
		return
	}
	//全てのユーザー情報を保存し、初期化する。
	aBot.guildCreateInitUsers(aSession, tGuildID)
	//全てのユーザーのステータスを保存し、初期化する。
	aBot.guildCreateInitStatuses(aSession, tGuildID, tMapRawData)
	//全てのロール情報を保存し、初期化する。
	aBot.guildCreateInitRoles(aSession, tGuildID)
}
func (aBot *Bot) guildCreateInitUsers(aSession *discordgo.Session, aGuildID string) {
	tMembers, tError := aSession.GuildMembers(aGuildID, "", 1000)
	if tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildCreate: failed to get members: %s", tError))
		return
	}
	for _, tMember := range tMembers {
		if tMember.User.Bot {
			continue
		}
		if tError := aBot.DataManager.updateUser(tMember, db.CurrentMember); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("guildCreate: failed to update user: %s", tError))
			return
		}
	}
}
func (aBot *Bot) guildCreateInitStatuses(aSession *discordgo.Session, aGuildID string, aMapRawData map[string]interface{}) {
	// これを実行中に onVoiceStateUpdate が起こると、死なないにしろ順序がおかしくなるかも。
	//全員分の音声状況を取得。初期状態にする。
	tInterfacePresences, tOk := aMapRawData["presences"].([]interface{})
	if !tOk {
		aBot.onEventError(aSession, "guildCreate: failed to convert presences")
		return
	}
	for _, tInterfacePresence := range tInterfacePresences {
		tPresence, tOk := tInterfacePresence.(map[string]interface{})
		if !tOk {
			aBot.onEventError(aSession, "guildCreate: failed to convert presence")
			continue
		}
		tInterfaceUser, tOk := tPresence["user"].(map[string]interface{})
		if !tOk {
			aBot.onEventError(aSession, "guildCreate: failed to convert member")
			continue
		}
		tUserID, tOk := tInterfaceUser["id"].(string)
		if !tOk {
			aBot.onEventError(aSession, "guildCreate: failed to convert user id")
			continue
		}

		// オフラインならスキップ
		if fmt.Sprint(tPresence["status"]) != fmt.Sprint(discordgo.StatusOnline) {
			continue
		}
		//Botはスキップ
		tMember, tError := aSession.GuildMember(aGuildID, tUserID)
		if tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("guildCreate: failed to get member: %s", tError))
			continue
		} else if tMember.User.Bot {
			continue
		}

		//ボイスチャンネルにいるならVoiceStateが手に入る。
		//いない場合ErrStateNotFoundがでるので、それを無視している。
		//ErrStateNotFoundが他の場合でもでるのなら、この処理はやめたほうがいいかもしれない。
		tVoiceState, tError := aSession.State.VoiceState(aGuildID, tUserID)
		if errors.Is(tError, discordgo.ErrStateNotFound) {
			tVoiceState = &discordgo.VoiceState{
				UserID:    tUserID,
				ChannelID: "",
			}
		} else if tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("guildCreate: failed to get voice state: %s", tError))
			continue
		}

		if tError := aBot.DataManager.updateStatus(tVoiceState, db.Online); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("guildCreate: failed to update status: %s", tError))
			return
		}
	}
}
func (aBot *Bot) guildCreateInitRoles(aSession *discordgo.Session, aGuildID string) {
	tRoles, tError := aSession.GuildRoles(aGuildID)
	if tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildCreate: failed to get roles: %s", tError))
		return
	}
	for _, tRole := range tRoles {
		aBot.DataManager.updateRole(tRole, db.CreateRole)
	}
}

// 誰かがサーバーに参加したとき
func (aBot *Bot) guildMemberAdd(aSession *discordgo.Session, aEvent *discordgo.GuildMemberAdd) {
	if aEvent.Member.User.Bot {
		return
	}
	if tError := aBot.DataManager.updateUser(aEvent.Member, db.CurrentMember); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildMemberAdd: failed to add user: %s", tError))
		return
	}
}

// 誰かのユーザー情報が更新されたとき
func (aBot *Bot) guildMemberUpdate(aSession *discordgo.Session, aEvent *discordgo.GuildMemberUpdate) {
	if aEvent.Member.User.Bot {
		return
	}
	if tError := aBot.DataManager.updateUser(aEvent.Member, db.CurrentMember); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildMemberUpdate: failed to update user: %s", tError))
		return
	}
}

// 誰かがサーバーから退出したとき
func (aBot *Bot) guildMemberRemove(aSession *discordgo.Session, aEvent *discordgo.GuildMemberRemove) {
	if aEvent.Member.User.Bot {
		return
	}
	if tError := aBot.DataManager.updateUser(aEvent.Member, db.OldMember); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildMemberRemove: failed to remove user: %s", tError))
		return
	}
}

// ロールが作成されたとき
func (aBot *Bot) roleCreate(aSession *discordgo.Session, aEvent *discordgo.GuildRoleCreate) {
	aBot.DataManager.updateRole(aEvent.Role, db.CreateRole)
}

// ロールが更新されたとき
func (aBot *Bot) roleUpdate(aSession *discordgo.Session, aEvent *discordgo.GuildRoleUpdate) {
	aBot.DataManager.updateRole(aEvent.Role, db.CreateRole)
}

// ロールが削除されたとき
func (aBot *Bot) roleDelete(aSession *discordgo.Session, aEvent *discordgo.GuildRoleDelete) {
	tRole := &discordgo.Role{
		ID: aEvent.RoleID,
	}
	aBot.DataManager.updateRole(tRole, db.DeleteRole)
}

// 誰かのオンライン状態が変わったとき
func (aBot *Bot) presenceUpdate(aSession *discordgo.Session, aEvent *discordgo.PresenceUpdate) {
	tVoiceState := &discordgo.VoiceState{
		UserID:    aEvent.User.ID,
		ChannelID: db.UnknownChannelID,
	}
	tIsOnline := db.Online
	if aEvent.Presence.Status != discordgo.StatusOnline {
		tIsOnline = db.Offline
	}

	if tError := aBot.DataManager.updateStatus(tVoiceState, tIsOnline); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("presenceUpdate: failed to update status: %s", tError))
		return
	}
}

// 誰かのボイスステータスが変わったとき。ボイスチャンネルに入ったり出たりしたとき。
func (aBot *Bot) voiceStateUpdate(aSession *discordgo.Session, aEvent *discordgo.VoiceStateUpdate) {
	if tError := aBot.DataManager.updateStatus(aEvent.VoiceState, db.UnknownOnline); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("voiceStateUpdate: failed to update status: %s", tError))
		return
	}
}

// 誰かがメッセージを送信したとき // apiの入出力
func (aBot *Bot) messageCreate(aSession *discordgo.Session, aEvent *discordgo.MessageCreate) {
	if aEvent.Author.ID == aSession.State.User.ID {
		return
	}

	// コマンド解析
	tContent := strings.Split(aEvent.Content, ",")
	if tContent[0] != PREFIX_COMMAND {
		return
	}
	if len(tContent) < 2 {
		if tError := aBot.sendMessageHelp(aSession, aEvent); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	}

	tCommand := tContent[1]

	if tCommand == "users" { //全てのユーザー情報を表示する。
		if tError := aBot.sendMessageAllUsers(aSession, aEvent); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	} else if tCommand == "status" { //指定したユーザーのステータスを表示する。
		tUsersID := tContent[5:]
		if tError := aBot.sendMessageStatuses(aSession, aEvent, tUsersID); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	} else if tCommand == "statuses" { //全てのユーザーのステータスを表示する。
		tUsersID, tError := api.GetAllUsersID(aBot.DataManager.DB)
		if tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to get all users id: %s", tError))
			return
		}
		if tError := aBot.sendMessageStatuses(aSession, aEvent, tUsersID); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	} else if tCommand == "roles" { //全てのロールを表示する。
		tRoles, tError := api.GetAllRoles(aBot.DataManager.DB)
		if tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to get all roles: %s", tError))
			return
		}
		if tError := aBot.sendMessageRoles(aSession, aEvent, tRoles); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	} else if tCommand == "statusesbyrole" { //ロールを指定して、そのロールを持っているユーザーを表示する。
		tUsersID, tError := api.GetUsersIDByRoleID(aBot.DataManager.DB, tContent[5])
		if tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to get users id by role id: %s", tError))
			return
		}
		if tError := aBot.sendMessageStatuses(aSession, aEvent, tUsersID); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	} else { //ヘルプを表示する。
		if tError := aBot.sendMessageHelp(aSession, aEvent); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	}
}

func (aBot *Bot) sendMessageAllUsers(aSession *discordgo.Session, aEvent *discordgo.MessageCreate) error {
	tUsers, tError := api.GetAllUsers(aBot.DataManager.DB)
	if tError != nil {
		return fmt.Errorf("failed to get all users: %s", tError)
	}
	tText := ""
	for _, tUser := range tUsers {
		tText += fmt.Sprintf("%v\n", *tUser)
		tText += "\n"
	}
	aSession.ChannelMessageSend(aEvent.ChannelID, tText)

	return nil
}
func (aBot *Bot) sendMessageStatuses(aSession *discordgo.Session, aEvent *discordgo.MessageCreate, aUsersID []string) error {
	// 引数パース
	tArgs := strings.Split(aEvent.Content, ",")
	tTimezoon, tError := time.LoadLocation("Asia/Tokyo")
	if tError != nil {
		return fmt.Errorf("failed to load location: %w", tError)
	}
	tStart, tError := time.ParseInLocation("20060102", tArgs[2], tTimezoon)
	if tError != nil {
		return fmt.Errorf("failed to parse start: %w", tError)
	}
	tEnd, tError := time.ParseInLocation("20060102", tArgs[3], tTimezoon)
	if tError != nil {
		return fmt.Errorf("failed to parse end: %w", tError)
	}
	tMinute, tError := strconv.ParseInt(tArgs[4], 10, 64)
	if tError != nil {
		return fmt.Errorf("failed to parse minute: %w", tError)
	}

	// 集計
	tStatuses, tError := api.GetTotalStatusesByUsersID(aBot.DataManager.DB, tStart, tEnd, time.Minute*time.Duration(tMinute), aUsersID)
	if tError != nil {
		return fmt.Errorf("failed to aggregate status within range by user ids: %w", tError)
	}

	// 結果表示
	aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("全集計範囲: %v/%v/%v - %v/%v/%v \n", tStatuses.Start.Year(), int(tStatuses.Start.Month()), tStatuses.Start.Day(),
		tStatuses.End.Year(), int(tStatuses.End.Month()), tStatuses.End.Day()))
	aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("Period: %v\n", tStatuses.Period))
	for tID, tStatuses := range tStatuses.StatusesByUserID {
		aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("UserID:%s\n", tID))
		for _, tStatus := range tStatuses {
			aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("集計範囲: %v - %v\n", tStatus.Start, tStatus.End))
			aSession.ChannelMessageSend(aEvent.ChannelID, "Channel\n")
			for tKey, tValue := range tStatus.ChannelByID {
				aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("%v : %v\n", tKey, tValue.TotalTime))
			}
			aSession.ChannelMessageSend(aEvent.ChannelID, "Voice\n")
			for tKey, tValue := range tStatus.VoiceByState {
				aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("%v : %v\n", tKey, tValue.TotalTime))
			}
			aSession.ChannelMessageSend(aEvent.ChannelID, "Online\n")
			for tKey, tValue := range tStatus.OnlineByStatus {
				aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("%v : %v\n", tKey, tValue.TotalTime))
			}
		}
	}
	return nil
}

func (aBot *Bot) sendMessageRoles(aSession *discordgo.Session, aEvent *discordgo.MessageCreate, aRoles []*db.Role) error {
	tText := ""
	for _, tRole := range aRoles {
		if tRole.Name == "@everyone" {
			continue
		}
		tText += fmt.Sprintf("%v\n", tRole)
		tText += "\n"
	}
	aSession.ChannelMessageSend(aEvent.ChannelID, tText)

	return nil
}

func (aBot *Bot) sendMessageHelp(aSession *discordgo.Session, aEvent *discordgo.MessageCreate) error {
	tMessage := PREFIX_COMMAND + ",users : 全てのユーザー情報\n" +
		PREFIX_COMMAND + ",statuses,20230101,20230131,1440 : 集計範囲と集計の分割時間(分)を指定しステータスを集計\n" +
		PREFIX_COMMAND + ",status,20230101,20230131,1440,USER_IDs... : 指定したUser(カンマ区切り)のステータスを集計\n" +
		PREFIX_COMMAND + ",statusesbyrole,20230101,20230131,1440,Role_ID : 指定したRoleを持つUserのステータスを集計\n" +
		PREFIX_COMMAND + ",roles : 全てのロール情報\n" +
		PREFIX_COMMAND + ",help : ヘルプを表示\n"

	aSession.ChannelMessageSend(aEvent.ChannelID, tMessage)
	return nil
}
