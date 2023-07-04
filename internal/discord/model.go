package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/one-scope/discord-time-tracker/internal/db"
	"github.com/robfig/cron/v3"
)

type OnlineStatus string

const (
	online        OnlineStatus = "online"
	offline       OnlineStatus = "offline"
	unknownOnline OnlineStatus = "unknown"
)

type dataManager struct {
	DataPathBase string
	UserByID     map[string]*db.User     // key: UserID
	StatusByID   map[string][]*statuslog // key: UserID
	DB           *db.SQLiteDB
}

// メモリ上のステータス構造体
type statuslog struct {
	UserID       string // キーがあるからいらないかも
	Timestamp    time.Time
	ChannelID    string
	VoiceState   string
	OnlineStatus OnlineStatus
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
