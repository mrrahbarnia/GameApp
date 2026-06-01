package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
}

type PostgreSQLDB struct {
	db *sql.DB
}

func New(cfg Config) *PostgreSQLDB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
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
