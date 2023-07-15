package db

import "time"

const (
	usersTable         = "users"
	usersTableID       = "discord_id"
	usersTableName     = "name"
	usersTableIsMember = "is_member"

	rolesTable   = "roles"
	rolesTableID = "role_id"

	statusesTable             = "statuses"
	statusesTableTimestamp    = "time_stamp"
	statusesTableChannelID    = "channel_id"
	statusesTableVoiceState   = "voice_state"
	statusesTableOnlineStatus = "online_status"
)

type IsMember bool

const (
	CurrentMember IsMember = true
	OldMember     IsMember = false
)

type OnlineStatus string

const (
	Online        OnlineStatus = "online"
	Offline       OnlineStatus = "offline"
	UnknownOnline OnlineStatus = "unknown"
)

type User struct {
	ID       string // DiscordのユーザーID
	Name     string
	IsMember IsMember //現在ギルドに所属しているか
	Roles    []string
}

type Statuslog struct {
	UserID       string
	Timestamp    time.Time
	ChannelID    string
	VoiceState   string
	OnlineStatus OnlineStatus
}
