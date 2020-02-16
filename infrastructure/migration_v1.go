package infrastructure

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var migrationV1 = gormigrate.Migration{
	ID: "1",
	Migrate: func(tx *gorm.DB) error {
		// Users
		if err := tx.AutoMigrate(&userModel{}).Error; err != nil {
			return err
		}

		// Documents
		if err := tx.AutoMigrate(&documentModel{}).Error; err != nil {
			return err
		}

		// Document Pages
		if err := tx.AutoMigrate(&documentPageModel{}).Error; err != nil {
			return err
		}

		return nil
	},

	Rollback: func(tx *gorm.DB) error {
		// Document Pages
		if err := tx.DropTable(documentPageModel{}.TableName()).Error; err != nil {
			return err
		}

		// Documents
		if err := tx.DropTable(documentModel{}.TableName()).Error; err != nil {
			return err
		}

		// Users
		if err := tx.DropTable(userModel{}.TableName()).Error; err != nil {
			return err
		}

		return nil
	},
}
