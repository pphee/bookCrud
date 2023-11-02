package pkg

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DbConnectGorm() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}
	return db
}
