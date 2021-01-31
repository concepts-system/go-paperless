package infrastructure

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrations = []*gormigrate.Migration{
	&migrationV1,
}

func buildMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}
