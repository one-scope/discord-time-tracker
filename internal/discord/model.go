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
	unknownOnline OnlineStatus = "unknown"
)

type dataManager struct {
	DataPathBase  string
	UsersMutex    sync.Mutex
	StatusesMutex sync.Mutex
	Users         map[string]*user        // key: UserID
	Statuses      map[string][]*statuslog // key: UserID
}

// メモリ上、ファイルの書き出し読み込み共用のユーザー構造体
type user struct {
	UserID   string   `json:"user_id"` // キーがあるからいらないかも
	UserName string   `json:"user_name"`
	IsMember IsMember `json:"is_member"` //現在ギルドに所属しているか
	Roles    []string `json:"roles"`
}

// メモリ上のステータス構造体
type statuslog struct {
	UserID       string // キーがあるからいらないかも
	Timestamp    time.Time
	ChannelID    string
	VoiceState   string
	OnlineStatus OnlineStatus
}

type channeldata struct {
	ChannelID string
	Time      int
}
type voicedata struct {
	VoiceState string
	Time       int
}

// データファイルの書き出し、読み込み用のステータス構造体
type statusdata struct {
	StartTime  time.Time // キーがあるからいらないかも
	StopTime   time.Time // キーがあるからいらないかも
	UserID     string    // キーがあるからいらないかも
	Channel    []*channeldata
	VoiceState []*voicedata
	OnlineTime int
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
