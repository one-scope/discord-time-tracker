package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/one-scope/discord-time-tracker/internal/config"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func New(aDBConfig *config.DBConfig) (*PostgresDB, error) {
	tDataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		aDBConfig.Host, aDBConfig.Port, aDBConfig.User, aDBConfig.Password, aDBConfig.DBName)
	tDb, tError := sqlx.Connect("postgres", tDataSourceName)
	if tError != nil {
		return nil, tError
	}

	return &PostgresDB{tDb}, nil
}

func (aDB *PostgresDB) Close() error {
	return aDB.DB.Close()
}
