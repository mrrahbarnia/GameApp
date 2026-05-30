package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgreSQLDB struct {
	db *sql.DB
}

func New() *PostgreSQLDB {
	connStr := "postgres://admin:123456@localhost:5432/db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Errorf("Cannot open PG connection due to: %w", err))
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &PostgreSQLDB{
		db: db,
	}
}
