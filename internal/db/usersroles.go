package db

import (
	"fmt"
)

// ロールテーブルの作成
func (aDB *PostgresDB) createUsersRolesTable() error {
	tQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT NOT NULL,%s TEXT NOT NULL,UNIQUE(%s,%s))",
		usersrolesTable, usersTableID, usersrolesTableID, usersTableID, usersrolesTableID)
	_, tError := aDB.DB.Exec(tQuery)
	return tError
}

// ユーザーIDを使ってmap[string]stringでロールを全て取得
func (aDB *PostgresDB) GetAllUsersRolesIDMapByUserID(aUserID string) (map[string]string, error) {
	tQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", usersrolesTableID, usersrolesTable, usersTableID)
	tRows, tError := aDB.DB.Query(tQuery, aUserID)
	if tError != nil {
		return nil, tError
	}
	defer tRows.Close()

	tRoles := map[string]string{}
	for tRows.Next() {
		tRole := ""
		if tError := tRows.Scan(&tRole); tError != nil {
			return nil, tError
		}
		tRoles[tRole] = tRole
	}
	//エラーチェック
	if tError := tRows.Err(); tError != nil {
		return nil, tError
	}
	return tRoles, nil
}

// ユーザーIDとロールIDのペアの存在確認
func (aDB *PostgresDB) IsExistsUserRoleByUserID(aUserID string, aRoleID string) (bool, error) {
	var tIsExists bool
	tQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1 AND %s = $2)", usersrolesTable, usersTableID, usersrolesTableID)
	if tError := aDB.DB.QueryRow(tQuery, aUserID, aRoleID).Scan(&tIsExists); tError != nil {
		return false, tError
	}
	return tIsExists, nil
}

// ユーザーIDを使ってロールを追加
func (aDB *PostgresDB) InsertUserRoleByUserID(aUserID string, aRoleID string) error {
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s) VALUES ($1,$2)", usersrolesTable, usersTableID, usersrolesTableID)
	_, tError := aDB.DB.Exec(tQuery, aUserID, aRoleID)
	return tError
}

// ユーザーIDを使ってロールを削除
func (aDB *PostgresDB) DeleteUserRoleByUserID(aUserID string, aRoleID string) error {
	tQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = $1 AND %s = $2", usersrolesTable, usersTableID, usersrolesTableID)
	_, tError := aDB.DB.Exec(tQuery, aUserID, aRoleID)
	return tError
}

// ユーザーIDを使ってロールを全て取得
func (aDB *PostgresDB) GetAllUsersRolesIDByUserID(aUserID string) ([]string, error) {
	tQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", usersrolesTableID, usersrolesTable, usersTableID)
	tRows, tError := aDB.DB.Query(tQuery, aUserID)
	if tError != nil {
		return nil, tError
	}
	defer tRows.Close()

	var tRoles []string
	for tRows.Next() {
		tRole := ""
		if tError := tRows.Scan(&tRole); tError != nil {
			return nil, tError
		}
		tRoles = append(tRoles, tRole)
	}
	//エラーチェック
	if tError := tRows.Err(); tError != nil {
		return nil, tError
	}
	return tRoles, nil
}

// 特定のロールを持つユーザーのIDを取得
func (aDB *PostgresDB) GetUsersIDByRoleID(aRoleID string) ([]string, error) {
	tQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", usersTableID, usersrolesTable, usersrolesTableID)
	tRows, tError := aDB.DB.Query(tQuery, aRoleID)
	if tError != nil {
		return nil, tError
	}
	defer tRows.Close()

	var tUsersID []string
	for tRows.Next() {
		tUserID := ""
		if tError := tRows.Scan(&tUserID); tError != nil {
			return nil, tError
		}
		tUsersID = append(tUsersID, tUserID)
	}
	//エラーチェック
	if tError := tRows.Err(); tError != nil {
		return nil, tError
	}
	return tUsersID, nil
}
