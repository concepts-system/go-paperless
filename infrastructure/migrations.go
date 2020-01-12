package infrastructure

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var migrations = []*gormigrate.Migration{
	&migrationV1,
}

func buildMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}
