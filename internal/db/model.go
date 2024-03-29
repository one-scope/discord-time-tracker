package db

import "time"

const (
	usersTable         = "users"
	usersTableID       = "discord_id"
	usersTableName     = "name"
	usersTableIsMember = "is_member"

	usersrolesTable   = "usersroles"
	usersrolesTableID = "usersrole_id"

	statusesTable             = "statuses"
	statusesTableID           = "statuses_id"
	statusesTableTimestamp    = "time_stamp"
	statusesTableChannelID    = "channel_id"
	statusesTableVoiceState   = "voice_state"
	statusesTableOnlineStatus = "online_status"

	rolesTable            = "roles"
	rolesTableID          = "roles_id"
	rolesTableName        = "roles_name"
	rolesTableManaged     = "roles_managed"
	rolesTableMentionable = "roles_mentionable"
	rolesTableHoist       = "roles_hoist"
	rolesTableColor       = "roles_color"
	rolesTablePosition    = "roles_position"
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

type VoiceState string

const (
	VoiceOffline VoiceState = "offline"
	VoiceDeaf    VoiceState = "speaker-mute"
	VoiceMute    VoiceState = "mic-mute"
	VoiceOn      VoiceState = "mic-on"
)

type RoleAction string

const (
	CreateRole RoleAction = "create"
	DeleteRole RoleAction = "delete"
)

type User struct {
	ID       string // DiscordのユーザーID
	Name     string
	IsMember IsMember //現在ギルドに所属しているか
	Roles    []string
}

type Statuslog struct {
	ID           string
	UserID       string
	Timestamp    time.Time
	ChannelID    string
	VoiceState   VoiceState
	OnlineStatus OnlineStatus
}

type Role struct {
	ID          string
	Name        string
	Managed     bool
	Mentionable bool
	Hoist       bool
	Color       int
	Position    int
}
