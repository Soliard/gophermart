package postgr

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (id, login, password_hash, created_at, roles, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Login, u.PasswordHash, u.CreatedAt, u.Roles, u.LastLoginAt)
	return err
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT * FROM users WHERE login = $1`
	err := r.db.GetContext(ctx, user, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UserExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE login = $1)`
	err := r.db.GetContext(ctx, &exists, query, login)
	return exists, err
}

func (r *UserRepository) UpdateLoginTime(ctx context.Context, userID string, time time.Time) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time, userID)
	return err
}
