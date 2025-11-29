package userrepository

import (
	"time"

	"github.com/codepnw/simple-bank/internal/features/user"
)

type UserModel struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func userDomainToModel(u *user.User) *UserModel {
	return &UserModel{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func userModelToDomain(u *UserModel) *user.User {
	return &user.User{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
