package discordbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

type DiscordBot struct {
	Session         *discordgo.Session
	Cron            *cron.Cron
	ExecutionTiming string
}

// 何かしらのイベントがあったとき。
func (aBot *DiscordBot) OnEvent(aHandler func(*discordgo.Session, *discordgo.Event)) {
	aBot.Session.AddHandler(aHandler)
}

func (aBot *DiscordBot) OnGuildCreate(aHandler func(*discordgo.Session, *discordgo.GuildCreate)) {
	aBot.Session.AddHandler(aHandler)
}
