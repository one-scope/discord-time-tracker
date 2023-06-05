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

	return nil
}

func (aBot *DiscordBot) setEventHandlers() {
	aBot.OnEvent(func(aSession *discordgo.Session, aEvent *discordgo.Event) {
		log.Println("discord event:", aEvent.Type) // デバッグ用にログ出力。
	})

}
