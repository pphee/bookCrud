package main

import (
	model "book/bookcrud"
	"book/server"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}

	db.AutoMigrate(&model.Book{})

	server.NewServer(db).StartGin()

}
