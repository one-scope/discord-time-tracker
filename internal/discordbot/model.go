package discordbot

import "github.com/bwmarrin/discordgo"

type DiscordBot struct {
	Session *discordgo.Session
}

// 何かしらのイベントがあったとき。
func (aBot *DiscordBot) OnEvent(aHandler func(*discordgo.Session, *discordgo.Event)) {
	aBot.Session.AddHandler(aHandler)
}
