package api

import "github.com/one-scope/discord-time-tracker/internal/db"

// 全てのロールを取得する
func GetAllRoles(aDB *db.PostgresDB) ([]*db.Role, error) {
	tRoles, tError := aDB.GetAllRoles()
	if tError != nil {
		return nil, tError
	}
	return tRoles, nil
}
