package db

type IsMember bool

const (
	CurrentMember IsMember = true
	OldMember     IsMember = false
)

// メモリ上、ファイルの書き出し読み込み共用のユーザー構造体
type User struct {
	ID       string
	Name     string
	IsMember IsMember //現在ギルドに所属しているか
	Roles    []string
}
