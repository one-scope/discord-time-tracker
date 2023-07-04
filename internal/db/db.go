package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	DB *sql.DB
}

const (
	userTable         = "users"
	userTableID       = "discord_id"
	userTableName     = "name"
	userTableIsMember = "is_member"
)

func New(aDBPath string) (*SQLiteDB, error) {
	if _, tError := os.Stat(aDBPath); os.IsNotExist(tError) {
		if tError := os.MkdirAll(filepath.Dir(aDBPath), 0755); tError != nil {
			return nil, tError
		}

		tFile, tError := os.Create(aDBPath)
		if tError != nil {
			return nil, tError
		}
		tFile.Close()
	}

	tDB, tError := sql.Open("sqlite3", aDBPath)
	if tError != nil {
		return nil, tError
	}

	if tError := tDB.Ping(); tError != nil {
		return nil, tError
	}

	tSQLiteDB := &SQLiteDB{tDB}

	if tError := tSQLiteDB.createUsersTable(); tError != nil {
		return nil, tError
	}
	if tError := tSQLiteDB.createUsersRolesTable(); tError != nil {
		return nil, tError
	}

	return tSQLiteDB, nil
}

func (aDB *SQLiteDB) Close() error {
	return aDB.DB.Close()
}

func (aDB *SQLiteDB) createUsersTable() error {
	_, tError := aDB.DB.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"%s INTEGER PRIMARY KEY NOT NULL,"+
			"%s TEXT NOT NULL,"+
			"%s BOOLEAN NOT NULL"+
			"created_at DATETIME DEFAULT (DATETIME('now','localtime')),"+
			"updated_at DATETIME DEFAULT (DATETIME('now','localtime'))"+
			")", userTable, userTableID, userTableName, userTableIsMember))
	return tError
}
func (aDB *SQLiteDB) IsExistsUserByID(aUserID string) (bool, error) {
	var tIsExists bool
	tError := aDB.DB.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = %s)", userTable, userTableID, aUserID)).Scan(&tIsExists)
	if tError != nil {
		return false, tError
	}
	return tIsExists, nil
}
func (aDB *SQLiteDB) createUsersRolesTable() error {
	_, tError := aDB.DB.Exec("CREATE TABLE IF NOT EXISTS roles (discord_id INTEGER NOT NULL, role INTEGER NOT NULL,created_at DATETIME DEFAULT (DATETIME('now','localtime')), updated_at DATETIME DEFAULT (DATETIME('now','localtime')), unique(discord_id, role))")
	return tError
}

// ユーザー新規登録
func (aDB *SQLiteDB) InsertUser(aDiscordID string, aName string, aIsMember IsMember) error {
	_, tError := aDB.DB.Exec(fmt.Sprintf("INSERT INTO %s (%s,%s,%s) VALUES (%s,%s,%t)", userTable, userTableID, userTableName, userTableIsMember, aDiscordID, aName, aIsMember))
	return tError
}

// ユーザー更新
func (aDB *SQLiteDB) UpdateUser(aDiscordID string, aName string, aIsMember IsMember) error {
	_, tError := aDB.DB.Exec(fmt.Sprintf("UPDATE %s SET (%s,%s) VALUES (%s,%t) WHERE %s = %s", userTable, userTableName, userTableIsMember, aName, aIsMember, userTableID, aDiscordID))
	return tError
}

func (aDB *SQLiteDB) GetAllUsers() (*[]*User, error) {
	tRows, tError := aDB.DB.Query(fmt.Sprintf("SELECT (%s,%s,%s) FROM %s", userTableID, userTableName, userTableIsMember, userTable))
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
	return &tUsers, nil
}
