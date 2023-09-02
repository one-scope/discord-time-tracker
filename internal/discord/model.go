package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/one-scope/discord-time-tracker/internal/db"
	"github.com/robfig/cron/v3"
)

type dataManager struct {
	DB                       *db.PostgresDB
	UsersByID                map[string]*db.User
	StatusesByID             map[string][]*db.Statuslog
	Roles                    []*roleWithAction
	PreViusStatusLogByUserID map[string]*db.Statuslog
}

type roleWithAction struct {
	Role   *db.Role
	Action db.RoleAction
}

type Bot struct {
	Session         *discordgo.Session
	Cron            *cron.Cron
	FlushTimingCron string
	ErrorChannel    string
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
func (aBot *Bot) onMessageCreate(aHandler func(*discordgo.Session, *discordgo.MessageCreate)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onGuildRoleCreate(aHandler func(*discordgo.Session, *discordgo.GuildRoleCreate)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onGuildRoleUpdate(aHandler func(*discordgo.Session, *discordgo.GuildRoleUpdate)) {
	aBot.Session.AddHandler(aHandler)
}
func (aBot *Bot) onGuildRoleDelete(aHandler func(*discordgo.Session, *discordgo.GuildRoleDelete)) {
	aBot.Session.AddHandler(aHandler)
}
