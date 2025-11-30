package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/codepnw/simple-bank/pkg/config"
	_ "github.com/lib/pq"
)

type DBExec interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

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

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) (err error)
}

type txManager struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) TxManager {
	return &txManager{db: db}
}

func (t *txManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Fatalf("rollback: %v", rbErr)
			}
		} else {
			cmErr := tx.Commit()
			err =cmErr
		}
	}()

	err = fn(tx)
	return err
}
