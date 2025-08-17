package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string     `json:"id" db:"id"`
	Login        string     `json:"login" db:"login"`
	PasswordHash string     `json:"-" db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	Roles        Roles      `json:"roles" db:"roles"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
}

func NewUser(login, passwordHash string) *User {
	return &User{
		ID:           uuid.New().String(),
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		Roles:        []Role{RoleUser},
	}
}
