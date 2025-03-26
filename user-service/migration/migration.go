package migration

import (
	"log"
	"userService/model"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {

	if db == nil {
		log.Fatal("Database connection is nil! Migration aborted.")
		return
	}

	err := db.AutoMigrate(
		&model.EmailVerification{},
		&model.PasswordReset{},
		&model.User{},
	)

	if err != nil {
		log.Fatal("Failed migration: ", err)
	}

	log.Println("Migration success!")
}
