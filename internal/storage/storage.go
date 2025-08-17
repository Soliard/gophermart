package storage

import (
	"context"

	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/repository"
	"github.com/Soliard/gophermart/internal/storage/postgr"
)

func New(ctx context.Context, c *config.Config) (repository.Storage, error) {
	return postgr.NewPostgresStorage(ctx, c.DatabaseDSN)
}
