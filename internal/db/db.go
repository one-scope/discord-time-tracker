package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sqlx.DB
}

// docker-composeの変数に依存している
func New() (*PostgresDB, error) {
	// 未実装：環境変数から引っ張ってdocker compose 変えてもいいようにする
	tDb, tError := sqlx.Connect("postgres", "host=postgres port=5432 user=postgres password=postgres dbname=db sslmode=disable")
	if tError != nil {
		return nil, tError
	}

	return &PostgresDB{tDb}, nil
}

func (aDB *PostgresDB) Close() error {
	return aDB.DB.Close()
}
