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
	tDb, tError := sqlx.Connect("postgres", tDataSourceName)
	if tError != nil {
		return nil, tError
	}

	tPostgresDB := &PostgresDB{tDb}

	// テーブル作成
	if tError := tPostgresDB.CreateUsersTable(); tError != nil {
		return nil, tError
	}
	if tError := tPostgresDB.CreateStatusesTable(); tError != nil {
		return nil, tError
	}
	if tError := tPostgresDB.CreateUsersRolesTable(); tError != nil {
		return nil, tError
	}
	if tError := tPostgresDB.CreateRolesTable(); tError != nil {
		return nil, tError
	}

	return tPostgresDB, nil
}

func (aDB *PostgresDB) Close() error {
	return aDB.DB.Close()
}
