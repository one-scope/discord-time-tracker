package api

import "github.com/one-scope/discord-time-tracker/internal/db"

// 全てのユーザーの情報を取得
func GetAllUsers(aDB *db.PostgresDB) ([]*db.User, error) {
	tUsers, tError := aDB.GetAllUsers()
	if tError != nil {
		return nil, tError
	}
	for _, tUser := range tUsers {
		tRoles, tError := aDB.GetAllRolesIDByUserID(tUser.ID)
		if tError != nil {
			return nil, tError
		}
		tUser.Roles = tRoles
	}

	return tUsers, nil
}

// 全てのユーザーのIDを取得
func GetAllUsersID(aDB *db.PostgresDB) ([]string, error) {
	tUsersID, tError := aDB.GetAllUsersID()
	if tError != nil {
		return nil, tError
	}

	return tUsersID, nil
}
