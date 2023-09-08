package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func New(aDBConfig *DBConfig) (*PostgresDB, error) {
	// DB接続
	tDataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		aDBConfig.Host, aDBConfig.Port, aDBConfig.User, aDBConfig.Password, aDBConfig.DBName)
	tDB, tError := sqlx.Connect("postgres", tDataSourceName)
	if tError != nil {
		return nil, tError
	}

	tPostgresDB := &PostgresDB{
		DB: tDB,
	}

	// テーブル作成
	if tError := tPostgresDB.createUsersTable(); tError != nil {
		return nil, tError
	}
	if tError := tPostgresDB.createStatusesTable(); tError != nil {
		return nil, tError
	}
	if tError := tPostgresDB.createUsersRolesTable(); tError != nil {
		return nil, tError
	}
	if tError := tPostgresDB.createRolesTable(); tError != nil {
		return nil, tError
	}

	return tPostgresDB, nil
}

func (aDB *PostgresDB) Close() error {
	return aDB.DB.Close()
}
