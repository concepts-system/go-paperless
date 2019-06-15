package migrations

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var migrations = []*gormigrate.Migration{
	&migrationV1,
}

// BuildMigrator builds a migrator for the database instance with the given migrations.
func BuildMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}
