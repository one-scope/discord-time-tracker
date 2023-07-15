package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// ユーザーテーブルを作成
func (aDB *PostgresDB) CreateUsersTable() error {
	tQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT PRIMARY KEY,%s TEXT NOT NULL,%s BOOLEAN NOT NULL)", usersTable, usersTableID, usersTableName, usersTableIsMember)
	_, tError := aDB.DB.Exec(tQuery)
	return tError
}

// ユーザーの存在確認
func (aDB *PostgresDB) IsExistsUserByID(aUserID string) (bool, error) {
	var tIsExists bool
	tQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", usersTable, usersTableID)
	if tError := aDB.DB.QueryRow(tQuery, aUserID).Scan(&tIsExists); tError != nil {
		return false, tError
	}
	return tIsExists, nil
}

// ユーザー新規登録
func (aDB *PostgresDB) InsertUser(aUser *User) error {
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s) VALUES ($1,$2,$3)", usersTable, usersTableID, usersTableName, usersTableIsMember)
	_, tError := aDB.DB.Exec(tQuery, aUser.ID, aUser.Name, aUser.IsMember)
	return tError
}

// ユーザー更新
func (aDB *PostgresDB) UpdateUser(aUser *User) error {
	tQuery := fmt.Sprintf("UPDATE %s SET %s = $1, %s = $2 WHERE %s = $3", usersTable, usersTableName, usersTableIsMember, usersTableID)
	_, tError := aDB.DB.Exec(tQuery, aUser.Name, aUser.IsMember, aUser.ID)
	return tError
}

// 全てのユーザーを取得
func (aDB *PostgresDB) GetAllUsers() ([]*User, error) {
	tQuery := fmt.Sprintf("SELECT * FROM %s", usersTable)
	tRows, tError := aDB.DB.Queryx(tQuery)
	if tError != nil {
		return nil, tError
	}
	defer tRows.Close()

	var tUsers []*User
	for tRows.Next() {
		tUser := User{}
		if tError := tRows.Scan(&tUser.ID, &tUser.Name, &tUser.IsMember); tError != nil {
			return nil, tError
		}
		tUsers = append(tUsers, &tUser)
	}
	//エラーチェック
	if tError := tRows.Err(); tError != nil {
		return nil, tError
	}
	return tUsers, nil
}
