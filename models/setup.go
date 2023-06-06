package models

import (
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open(os.Getenv("DB_PATH")), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Site{})
	if err != nil {
		return
	}

	DB = database
}
