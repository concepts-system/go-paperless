package infrastructure

import (
	"github.com/concepts-system/go-paperless/config"
	"github.com/concepts-system/go-paperless/errors"
	gorm_logrus "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database defines a struct holding all required information for accessing
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

	var dialector gorm.Dialector
	switch db.config.Database.Type {
	case "sqlite3":
		dialector = sqlite.Open(db.config.Database.URL)
	case "postgres":
		dialector = postgres.Open(db.config.Database.URL)
	default:
		log.Fatalf("Unknown database type '%s'", db.config.Database.Type)
	}

	db.DB, err = gorm.Open(
		dialector,
		&gorm.Config{
			Logger: gorm_logrus.New(),
		},
	)

	if err != nil {
		return errors.Wrapf(err, "Failed to connect to database")
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
