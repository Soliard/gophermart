package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Role string

type User struct {
	ID           string     `json:"id"`
	Login        string     `json:"login"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	Roles        []Role     `json:"roles"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}

func NewUser(login, passwordHash string) *User {
	return &User{
		ID:           uuid.New().String(),
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		Roles:        []Role{RoleUser},
	}
}
