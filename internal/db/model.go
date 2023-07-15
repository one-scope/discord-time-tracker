package db

import "time"

// テーブル名、カラム名
// 未実装：init.sqlと揃える必要がある。これ改善。goでテーブル作る
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

// メモリ上、ファイルの書き出し読み込み共用のユーザー構造体
// 未実装：Userモデルのdbとusers.goを揃える必要がある
type User struct {
	ID       string   `db:"discord_id"` // DiscordのユーザーID
	Name     string   `db:"name"`
	IsMember IsMember `db:"is_member"` //現在ギルドに所属しているか
	Roles    []string `db:"roles"`
}

// メモリ上、ファイルの書き出し読み込み共用のステータス構造体
// 未実装：Statuslogモデルのdbとstatuses.goを揃える必要がある
type Statuslog struct {
	UserID       string       `db:"discord_id"`
	Timestamp    time.Time    `db:"time_stamp"`
	ChannelID    string       `db:"channel_id"`
	VoiceState   string       `db:"voice_state"`
	OnlineStatus OnlineStatus `db:"online_status"`
}

//疑問：ここら辺の`db:"hogehoge"`やめた方がいいかな、:discord_idじゃなくて$1にそろえた方がいいかな？
