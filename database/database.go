package database

import (
	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/errors"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var database *gorm.DB

// DB returns the current database/ORM instance used in the persistence context
// of the application.
func DB() *gorm.DB {
	if database == nil {
		var err error
		database, err = gorm.Open(common.Config().GetDatabaseType(), common.Config().GetDatabaseURL().String())

		if err != nil {
			panic(errors.Wrapf(err, "Failed to connect to database"))
		}

		// Run auto-migration for development environments
		if common.Config().IsDevelopment() {
			database.LogMode(true)
		}
	}

	return database
}

// BuildMigrator builds a migrator for the database instance with the given migrations.
func BuildMigrator(migrations []*gormigrate.Migration) *gormigrate.Gormigrate {
	return gormigrate.New(database, gormigrate.DefaultOptions, migrations)
}
