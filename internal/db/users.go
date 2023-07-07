package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s) VALUES (:discord_id,:name,:is_member)", usersTable, usersTableID, usersTableName, usersTableIsMember)
	_, tError := aDB.DB.NamedExec(tQuery, *aUser)
	return tError
}

// ユーザー更新
func (aDB *PostgresDB) UpdateUser(aUser *User) error {
	tQuery := fmt.Sprintf("UPDATE %s SET %s = :name, %s = :is_member WHERE %s = :discord_id", usersTable, usersTableName, usersTableIsMember, usersTableID)
	_, tError := aDB.DB.NamedExec(tQuery, *aUser)
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
		if tError := tRows.StructScan(&tUser); tError != nil {
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
