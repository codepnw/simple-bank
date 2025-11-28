package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/codepnw/simple-bank/pkg/config"
)

func ConnectPostgres(cfg *config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db failed: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db failed: %w", err)
	}
	return db, err
}
