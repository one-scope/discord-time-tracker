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

func New(aConfig *config.DiscordBotConfig, aDB *db.PostgresDB) (*Bot, error) {
	tSession, tError := discordgo.New("Bot " + aConfig.DiscordBotToken)
	if tError != nil {
		return nil, tError
	}
	tSession.Identify.Intents = discordgo.IntentsAll // 現在、テストのため全て許可
	tCron := cron.New()
	tManager := &dataManager{
		UsersByID:                      map[string]*db.User{},
		StatusesByID:                   map[string][]*db.Statuslog{},
		DB:                             aDB,
		PreViusStatusLogIDByUserID:     map[string]string{},
		PreViusStatusLogOnlineByUserID: map[string]db.OnlineStatus{},
	}
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
	// イベント実行
	aBot.setEventHandlers()
	if tError := aBot.Session.Open(); tError != nil {
		return tError
	}

	// 定期実行
	if _, tError := aBot.Cron.AddFunc(aBot.FlushTimingCron, aBot.DataManager.flushData()); tError != nil {
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
		log.Println(aErrorMessage)
		return
	}
	if _, tError := aSession.ChannelMessageSend(aBot.ErrorChannel, aErrorMessage); tError != nil {
		log.Println("Error: failed to send message to error channel")
		log.Println(tError)
		return
	}
}

// 未実装：エラー時に指定したチャンネルにエラーを送信し、続行。
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
}
func (aBot *Bot) event(aSession *discordgo.Session, aEvent *discordgo.Event) {
	if aEvent.Type == "GUILD_CREATE" {
		aBot.guildCreate(aSession, aEvent)
	}
}
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
func (aBot *Bot) guildMemberAdd(aSession *discordgo.Session, aEvent *discordgo.GuildMemberAdd) {
	if aEvent.Member.User.Bot {
		return
	}
	if tError := aBot.DataManager.updateUser(aEvent.Member, db.CurrentMember); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildMemberAdd: failed to add user: %s", tError))
		return
	}
}
func (aBot *Bot) guildMemberUpdate(aSession *discordgo.Session, aEvent *discordgo.GuildMemberUpdate) {
	if aEvent.Member.User.Bot {
		return
	}
	if tError := aBot.DataManager.updateUser(aEvent.Member, db.CurrentMember); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildMemberUpdate: failed to update user: %s", tError))
		return
	}
}
func (aBot *Bot) guildMemberRemove(aSession *discordgo.Session, aEvent *discordgo.GuildMemberRemove) {
	if aEvent.Member.User.Bot {
		return
	}
	if tError := aBot.DataManager.updateUser(aEvent.Member, db.OldMember); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("guildMemberRemove: failed to remove user: %s", tError))
		return
	}
}
func (aBot *Bot) presenceUpdate(aSession *discordgo.Session, aEvent *discordgo.PresenceUpdate) {
	tVoiceState := &discordgo.VoiceState{
		UserID:    aEvent.User.ID,
		ChannelID: "",
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
func (aBot *Bot) voiceStateUpdate(aSession *discordgo.Session, aEvent *discordgo.VoiceStateUpdate) {
	if tError := aBot.DataManager.updateStatus(aEvent.VoiceState, db.UnknownOnline); tError != nil {
		aBot.onEventError(aSession, fmt.Sprintf("voiceStateUpdate: failed to update status: %s", tError))
		return
	}
}

// 集計APIのテスト用
func (aBot *Bot) messageCreate(aSession *discordgo.Session, aEvent *discordgo.MessageCreate) {
	if aEvent.Author.ID == aSession.State.User.ID {
		return
	}

	if aEvent.Content == "allusers" {
		if tError := aBot.sendMessageAllUsers(aSession, aEvent); tError != nil {
			aBot.onEventError(aSession, fmt.Sprintf("messageCreate: failed to send message: %s", tError))
		}
	} else {
		if tError := aBot.sendMessageAllStatuses(aSession, aEvent); tError != nil {
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
func (aBot *Bot) sendMessageAllStatuses(aSession *discordgo.Session, aEvent *discordgo.MessageCreate) error {
	// status,20230701,20230731,60,684946062061993994
	// status, start , end , minute , userIDs...

	// 引数パース
	tArgs := strings.Split(aEvent.Content, ",")
	tStart, tError := time.Parse("20060102", tArgs[1])
	if tError != nil {
		return fmt.Errorf("failed to parse start: %w", tError)
	}
	tEnd, tError := time.Parse("20060102", tArgs[2])
	if tError != nil {
		return fmt.Errorf("failed to parse end: %w", tError)
	}
	tMinute, tError := strconv.ParseInt(tArgs[3], 10, 64)
	if tError != nil {
		return fmt.Errorf("failed to parse minute: %w", tError)
	}

	// 集計
	tStatuses, tError := api.AggregateStatusWithinRangeByUserIDs(aBot.DataManager.DB, tStart, tEnd, time.Minute*time.Duration(tMinute), tArgs[4:])
	if tError != nil {
		return fmt.Errorf("failed to aggregate status within range by user ids: %w", tError)
	}

	// 結果表示
	aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("Start: %v\n\n", tStatuses.Start))
	aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("End: %v\n\n", tStatuses.End))
	aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("Period: %v\n\n", tStatuses.Period))
	aSession.ChannelMessageSend(aEvent.ChannelID, "StatusesByUserID:\n\n")
	for tID, tStatuses := range tStatuses.StatusesByUserID {
		aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("USER:%s\n\n", tID))
		for _, tStatus := range tStatuses {
			aSession.ChannelMessageSend(aEvent.ChannelID, "Start-End\n\n")
			aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("%v-%v\n\n", tStatus.Start, tStatus.End))
			aSession.ChannelMessageSend(aEvent.ChannelID, "Channel\n\n")
			aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("%v\n\n", tStatus.ChannelByID))
			aSession.ChannelMessageSend(aEvent.ChannelID, "Voice\n\n")
			aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("%v\n\n", tStatus.VoiceByState))
			aSession.ChannelMessageSend(aEvent.ChannelID, "Online\n\n")
			aSession.ChannelMessageSend(aEvent.ChannelID, fmt.Sprintf("%v\n\n", tStatus.OnlineByStatus))
		}
	}
	return nil
}
