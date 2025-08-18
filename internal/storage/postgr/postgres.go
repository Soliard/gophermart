package postgr

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

func NewConnection(ctx context.Context, connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connectionString)
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

	return db, nil
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

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
