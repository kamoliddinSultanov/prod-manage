package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(file, dsn string) error {
	m, err := migrate.New(file, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}
