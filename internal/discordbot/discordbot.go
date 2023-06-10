package discordbot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func New(aToken string) (*DiscordBot, error) {
	tSession, tError := discordgo.New(aToken)
	if tError != nil {
		return nil, tError
	}
	tBot := &DiscordBot{
		Session: tSession,
	}
	return tBot, nil
}

func Start(aBot *DiscordBot) error {
	aBot.setEventHandlers()
	if tError := aBot.Session.Open(); tError != nil {
		return tError
	}
	// tTicker := time.NewTicker(5 * time.Minute)
	// go func() {
	// 	for range tTicker.C {
	// 		log.Println("flush discord activity")
	// 		if tError := tSession.flushStatuses(); tError != nil {
	// 			log.Println("failed to flush statuses:", tError)
	// 		}
	// 	}
	// }()

	return nil
}

func (aBot *DiscordBot) setEventHandlers() {
	aBot.OnEvent(func(aSession *discordgo.Session, aEvent *discordgo.Event) {
		log.Println("discord event:", aEvent.Type) // デバッグ用にログ出力。
	})
	// サーバーに接続したとき。
	aBot.OnGuildCreate(func(aSession *discordgo.Session, aEvent *discordgo.GuildCreate) {
		// 全員分の音声状況を取得。初期状態にする。
		// これを実行中に OnVoiceStateUpdate が起こると、死なないにしろ順序がおかしくなるかも。
		for _, tMember := range aEvent.Members {
			tVoiceState, tError := aSession.State.VoiceState(tMember.GuildID, tMember.User.ID)
			if tError != nil {
				log.Println("failed to get voice state:", tError)
				continue
			}
			log.Println("Voice State:", tVoiceState)

			// if tError := aDiscordActivity.updateStatus(aSession, tVoiceState); tError != nil {
			// 	log.Println("failed to update status:", tError)
			// }
		}
	})

	// // 誰かの音声通話が更新されたとき。
	// // 接続、切断もこれ。切断時は ChannelID が空文字。
	// aBot.OnVoiceStateUpdate(func(aSession *discordgo.Session, aEvent *discordgo.VoiceStateUpdate) {
	// 	if tError := aDiscordActivity.updateStatus(aSession, aEvent.VoiceState); tError != nil {
	// 		log.Println("failed to update status:", tError)
	// 	}
	// })
}
