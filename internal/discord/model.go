package discord

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

type IsMember bool
type OnlineStatus string

const (
	currentMember IsMember     = true
	oldMember     IsMember     = false
	online        OnlineStatus = "online"
	offline       OnlineStatus = "offline"
	unknownOnline OnlineStatus = "unknown" //
)

type dataManager struct {
	DataDir       string
	UsersMutex    sync.Mutex
	StatusesMutex sync.Mutex
	Users         map[string]*user     // key: UserID
	Statuses      map[string][]*status // key: UserID
}

type user struct {
	UserID   string
	UserName string
	IsMember IsMember //現在ギルドに所属しているか
	Roles    []string
}

type status struct {
	UserID       string
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

// 何かしらのイベントがあったとき。
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
