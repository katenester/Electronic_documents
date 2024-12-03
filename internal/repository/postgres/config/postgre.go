package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	UsersTable          = "users"
	DocumentsTable      = "documents"
	DocumentGrantsTable = "document_grants"
	SessionsTable       = "sessions"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	// Check connect
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}
