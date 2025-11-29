package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/codepnw/simple-bank/internal/features/user"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
)
//go:generate mockgen -source=user_repository.go -destination=mock_user_repository.go -package=userrepository
type UserRepository interface {
	Insert(ctx context.Context, input *user.User) (*user.User, error)
	FindByID(ctx context.Context, id int64) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Insert(ctx context.Context, input *user.User) (*user.User, error) {
	m := userDomainToModel(input)
	query := `
		INSERT INTO users (username, password, first_name, last_name, email)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		m.Username,
		m.Password,
		m.FirstName,
		m.LastName,
		m.Email,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_username_key"`) {
			return nil, errs.ErrUsernameAlreadyExists
		}
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
			return nil, errs.ErrEmailAlreadyExists
		}
		return nil, err
	}
	return userModelToDomain(m), nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email FROM users
		WHERE id = $1 LIMIT 1
	`
	u := new(user.User)
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.Email,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, password FROM users
		WHERE email = $1 LIMIT 1
	`
	u := new(user.User)
	if err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}
