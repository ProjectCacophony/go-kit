package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	// we use the file source to load the migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run performs migrations all the way up
func Run(logger *zap.Logger, path, dsn string) error {
	m, err := migrate.New(
		"file://"+path,
		dsn,
	)
	if err != nil {
		return err
	}
	defer m.Close()
	m.Log = NewLogger(logger, false)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	migrateVersion, migrateDirty, err := m.Version()
	if err != nil {
		return err
	}

	logger.Info(
		"migrations completed successfully",
		zap.Uint("version", migrateVersion),
		zap.Bool("dirty", migrateDirty),
	)

	return nil
}
