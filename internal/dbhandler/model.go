package dbhandler

type IsMember bool

const (
	CurrentMember IsMember = true
	OldMember     IsMember = false
)

// メモリ上、ファイルの書き出し読み込み共用のユーザー構造体
type User struct {
	UserID   string   `json:"user_id"` // キーがあるからいらないかも
	UserName string   `json:"user_name"`
	IsMember IsMember `json:"is_member"` //現在ギルドに所属しているか
	Roles    []string `json:"roles"`
}
