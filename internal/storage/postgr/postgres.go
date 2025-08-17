package postgr

import (
	"context"
	"embed"

	"github.com/Soliard/gophermart/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	userRepository  repository.UserRepositoryInterface
	orderRepository repository.OrderRepositoryInterface
}

func (s *PostgresStorage) UserRepository() repository.UserRepositoryInterface {
	return s.userRepository
}

func (s *PostgresStorage) OrderRepository() repository.OrderRepositoryInterface {
	return s.orderRepository
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewPostgresStorage(ctx context.Context, databaseDSN string) (repository.Storage, error) {
	db, err := sqlx.Open("postgres", databaseDSN)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		userRepository:  newUserRepository(db),
		orderRepository: newOrderRepository(db),
	}, nil
}

func runMigrations(db *sqlx.DB) error {
	//exmpl migrate create -ext sql -dir migrations -seq create_metrics_table
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	migr, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return err
	}

	if err := migr.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
