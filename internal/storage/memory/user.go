package memory

import (
	"context"
	"strings"
	"time"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
)

type Users map[string]*models.User

type UserRepository struct {
	Users Users
}

func NewUserRepository() *UserRepository {
	users := Users{}
	return &UserRepository{Users: users}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	u.CreatedAt = time.Now()
	r.Users[strings.ToLower(u.Login)] = u
	return nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	login = strings.ToLower(login)
	if u, ok := r.Users[login]; ok {
		return u, nil
	}
	return nil, errs.UserNotFound
}

func (r *UserRepository) UserExists(ctx context.Context, login string) (bool, error) {
	login = strings.ToLower(login)
	_, ok := r.Users[login]
	return ok, nil
}
