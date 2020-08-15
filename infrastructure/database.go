package infrastructure

import (
	"github.com/concepts-system/go-paperless/config"
	"github.com/concepts-system/go-paperless/errors"
	log "github.com/kpango/glg"

	"github.com/jinzhu/gorm"
)

// Database defines a struct holding all required infromation for accessing
// the database.
type Database struct {
	*gorm.DB
	config *config.Configuration
}

// NewDatabase creates a new database using the given configuration.
func NewDatabase(config *config.Configuration) *Database {
	if config == nil {
		log.Fatal("Configuration may not be null!")
	}

	return &Database{
		config: config,
	}
}

// Connect tries to establish a connection to the configured database.
func (db *Database) Connect() error {
	var err error
	db.DB, err = gorm.Open(
		db.config.Database.Type,
		db.config.Database.URL,
	)

	if err != nil {
		return errors.Wrapf(err, "Failed to connect to database")
	}

	if !db.config.IsProductionMode() {
		db.DB.LogMode(true)
	}

	return nil
}

// Migrate migrates the given database instance to the latest version.
func (db *Database) Migrate() error {
	return buildMigrator(db.DB).Migrate()
}

// MigrateTo migrates the given database instance to the given version.
func (db *Database) MigrateTo(version string) error {
	return buildMigrator(db.DB).MigrateTo(version)
}
