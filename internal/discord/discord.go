package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/one-scope/discord-time-tracker/internal/config"
	"github.com/robfig/cron/v3"
)

func New(aConfig *config.DiscordBotConfig) (*Bot, error) {
	tSession, tError := discordgo.New("Bot " + aConfig.DiscordBotToken)
	if tError != nil {
		return nil, tError
	}
	tSession.Identify.Intents = discordgo.IntentsAll // 現在、テストのため全て許可
	tCron := cron.New()
	tManager := &dataManager{
		DataPathBase:  aConfig.DataPathBase,
		UsersMutex:    sync.Mutex{},
		StatusesMutex: sync.Mutex{},
		Users:         map[string]*user{},
		Statuses:      map[string][]*statuslog{},
	}
	tBot := &Bot{
		Session:         tSession,
		Cron:            tCron,
		ExecutionTiming: aConfig.ExecutionTiming,
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
	if _, tError := aBot.Cron.AddFunc(aBot.ExecutionTiming, aBot.DataManager.flushData()); tError != nil {
		return tError
	}
	aBot.Cron.Start()

	return nil
}

func (aBot *Bot) Close() error {
	aBot.Cron.Stop()
	return aBot.Session.Close()
}

// エラー時に指定したチャンネルにエラーを送信し、続行。
func (aBot *Bot) setEventHandlers() {
	// ボットが起動したとき。
	// 未実装：ボットが起動したときに、ユーザーの情報(ID,Name,Role等)を取得する。
	aBot.onEvent(func(aSession *discordgo.Session, aEvent *discordgo.Event) {
		var tMapRawData map[string]interface{}
		json.Unmarshal(aEvent.RawData, &tMapRawData)

		if aEvent.Type == "GUILD_CREATE" {
			aBot.guildCreate(aSession, aEvent, tMapRawData)
		}
	})
	// 誰かがサーバーに参加したとき。
	aBot.onGuildMemberAdd(func(aSession *discordgo.Session, aEvent *discordgo.GuildMemberAdd) {
		if tError := aBot.DataManager.updateUser(aEvent.Member, currentMember); tError != nil {
			log.Println("failed to add user:", tError)
		}
	})
	//誰かのロールが変わったとき。
	aBot.onGuildMemberUpdate(func(aSession *discordgo.Session, aEvent *discordgo.GuildMemberUpdate) {
		if tError := aBot.DataManager.updateUser(aEvent.Member, currentMember); tError != nil {
			log.Println("failed to update user:", tError)
		}
	})
	// 誰かがサーバーから退出したとき。
	aBot.onGuildMemberRemove(func(aSession *discordgo.Session, aEvent *discordgo.GuildMemberRemove) {
		if tError := aBot.DataManager.updateUser(aEvent.Member, oldMember); tError != nil {
			log.Println("failed to remove user:", tError)
		}
	})
	//誰かのオンラインになったとき
	aBot.onPresenceUpdate(func(aSession *discordgo.Session, aEvent *discordgo.PresenceUpdate) {
		tVoiceState := &discordgo.VoiceState{
			UserID:    aEvent.User.ID,
			ChannelID: "",
		}
		tIsOnline := online
		if aEvent.Presence.Status != discordgo.StatusOnline {
			tIsOnline = offline
		}

		if tError := aBot.DataManager.updateStatus(tVoiceState, tIsOnline); tError != nil {
			log.Println("failed to update status:", tError)
		}
	})
	// 誰かの音声通話が更新されたとき。// 接続、切断もこれ。切断時は ChannelID が空文字。
	aBot.onVoiceStateUpdate(func(aSession *discordgo.Session, aEvent *discordgo.VoiceStateUpdate) {
		if tError := aBot.DataManager.updateStatus(aEvent.VoiceState, unknownOnline); tError != nil {
			log.Println("failed to update status:", tError)
		}
	})
}

// ステータスだけじゃなくて、ユーザー情報も取得する。
func (aBot *Bot) guildCreate(aSession *discordgo.Session, aEvent *discordgo.Event, aRawData map[string]interface{}) {
	//全員分の音声状況を取得。初期状態にする。
	// これを実行中に onVoiceStateUpdate が起こると、死なないにしろ順序がおかしくなるかも。

	//GuildIDを取得
	tGuildID, tOk := aRawData["id"].(string)
	if !tOk {
		log.Println("failed to convert guild id")
		return
	}
	tInterfacePresences, tOk := aRawData["presences"].([]interface{})
	if !tOk {
		log.Println("failed to convert presences")
		return
	}
	for _, tInterfacePresence := range tInterfacePresences {
		tPresence, tOk := tInterfacePresence.(map[string]interface{})
		if !tOk {
			log.Println("failed to convert presence")
			continue
		}
		// オフラインならスキップ
		if fmt.Sprint(tPresence["status"]) != fmt.Sprint(discordgo.StatusOnline) {
			continue
		}

		//UserIDを取得
		tInterfaceUser, tOk := tPresence["user"].(map[string]interface{})
		if !tOk {
			log.Println("failed to convert member")
			continue
		}
		tUserID := tInterfaceUser["id"].(string)

		//Botはスキップ
		tMember, tError := aSession.GuildMember(tGuildID, tUserID)
		if tError != nil {
			log.Println("failed to get member:", tError)
			continue
		} else if tMember.User.Bot {
			continue
		}

		//ボイスチャンネルにいるならVoiceStateが手に入る。
		//いない場合ErrStateNotFoundがでるので、それを無視している。
		//ErrStateNotFoundが他の原因ででるかは不明。
		tVoiceState, tError := aSession.State.VoiceState(tGuildID, tUserID)
		if errors.Is(tError, discordgo.ErrStateNotFound) {
			tVoiceState = &discordgo.VoiceState{
				UserID:    tUserID,
				ChannelID: "",
			}
		} else if tError != nil {
			log.Println("failed to get voice state:", tError)
			continue
		}

		if tError := aBot.DataManager.updateStatus(tVoiceState, online); tError != nil {
			log.Println("failed to update status:", tError)
		}
	}
}
