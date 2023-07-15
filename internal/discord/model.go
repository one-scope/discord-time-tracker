package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/one-scope/discord-time-tracker/internal/db"
	"github.com/robfig/cron/v3"
)

type dataManager struct {
	UsersByID    map[string]*db.User        // key: UserID
	StatusesByID map[string][]*db.Statuslog // key: UserID
	DB           *db.PostgresDB
}

type Bot struct {
	Session         *discordgo.Session
	Cron            *cron.Cron
	ExecutionTiming string
	DataManager     *dataManager
}

func (aBot *Bot) onEvent(aHandler func(*discordgo.Session, *discordgo.Event)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onVoiceStateUpdate(aHandler func(*discordgo.Session, *discordgo.VoiceStateUpdate)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onGuildMemberAdd(aHandler func(*discordgo.Session, *discordgo.GuildMemberAdd)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onGuildMemberUpdate(aHandler func(*discordgo.Session, *discordgo.GuildMemberUpdate)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onGuildMemberRemove(aHandler func(*discordgo.Session, *discordgo.GuildMemberRemove)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onPresenceUpdate(aHandler func(*discordgo.Session, *discordgo.PresenceUpdate)) {
	aBot.Session.AddHandler(aHandler)
}
