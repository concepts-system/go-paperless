package migrations

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"

	"github.com/concepts-system/go-paperless/users"
	"github.com/concepts-system/go-paperless/documents"
)

var migrationV1 = gormigrate.Migration{
	ID: "1",
	Migrate: func(tx *gorm.DB) error {
		// Users
		if err := tx.AutoMigrate(&users.UserModel{}).Error; err != nil {
			return err
		}

		// Documents
		if err := tx.AutoMigrate(&documents.DocumentModel{}).Error; err != nil {
			return err
		}

		// Document Pages
		if err := tx.AutoMigrate(&documents.PageModel{}).Error; err != nil {
			return err
		}

		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		// Document Pages
		if err := tx.DropTable(documents.PageModel{}.TableName()).Error; err != nil {
			return err
		}

		// Documents
		if err := tx.DropTable(documents.DocumentModel{}.TableName()).Error; err != nil {
			return err
		}

		// Users
		if err := tx.DropTable(users.UserModel{}.TableName()).Error; err != nil {
			return err
		}

		return nil
	},
}
