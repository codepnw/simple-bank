package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

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

	counts := 0
	for {
		err := db.Ping()
		if err == nil {
			log.Println("connect database successfully")
			break
		}

		if counts >= 15 {
			// > 30 sec (15 * 2)
			return nil, fmt.Errorf("failed to connect database after retries: %w", err)
		}

		log.Printf("database not ready retrying in 2 seconds (%d/15)", counts+1)
		time.Sleep(2 * time.Second)
		counts++
	}

	return db, err
}

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) (err error)
}

type txManager struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) (TxManager, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &txManager{db: db}, nil
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
