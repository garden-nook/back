package migrator

import (
	"errors"
	"fmt"
	"garden-nook/migrations"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Run применяет все pending-миграции к базе данных.
// Возвращает nil, если миграции успешно применены или уже актуальны.
func Run(dsn string, log *slog.Logger) error {
	log.Info("starting database migration")

	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	dbDriver, err := database.Open(dsn)
	if err != nil {
		return fmt.Errorf("create db driver: %w", err)
	}
	defer dbDriver.Close()

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			version, _, _ := m.Version()
			log.Info("database is up to date", "version", version)
			return nil
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	version, _, _ := m.Version()
	log.Info("migrations applied successfully", "version", version)
	return nil
}
