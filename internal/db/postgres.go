package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func MustConnectPostgres(pgConfig *PostgresConfig) *sqlx.DB {
	db := sqlx.MustConnect("postgres", pgConfig.String())
	return db
}
