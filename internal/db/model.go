package db

type IsMember bool

// テーブル名、カラム名
// Userモデルとも依存している
const (
	usersTable         = "users"
	usersTableID       = "discord_id"
	usersTableName     = "name"
	usersTableIsMember = "is_member"

	rolesTable   = "roles"
	rolesTableID = "role_id"
)

const (
	CurrentMember IsMember = true
	OldMember     IsMember = false
)

// メモリ上、ファイルの書き出し読み込み共用のユーザー構造体
type User struct {
	ID       string   `db:"discord_id"` // DiscordのユーザーID
	Name     string   `db:"name"`
	IsMember IsMember `db:"is_member"` //現在ギルドに所属しているか
	Roles    []string `db:"roles"`
}
