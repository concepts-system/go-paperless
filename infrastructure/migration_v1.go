package infrastructure

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrationV1 = gormigrate.Migration{
	ID: "1",
	Migrate: func(tx *gorm.DB) error {
		// Users
		if err := tx.AutoMigrate(&userModel{}); err != nil {
			return err
		}

		// Documents
		if err := tx.AutoMigrate(&documentModel{}); err != nil {
			return err
		}

		// Document Pages
		if err := tx.AutoMigrate(&documentPageModel{}); err != nil {
			return err
		}

		return nil
	},

	Rollback: func(tx *gorm.DB) error {
		// Document Pages
		if err := tx.Migrator().DropTable(documentPageModel{}.TableName()); err != nil {
			return err
		}

		// Documents
		if err := tx.Migrator().DropTable(documentModel{}.TableName()); err != nil {
			return err
		}

		// Users
		if err := tx.Migrator().DropTable(userModel{}.TableName()); err != nil {
			return err
		}

		return nil
	},
}
