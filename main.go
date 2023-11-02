package main

import (
	model "book/bookcrud"
	"book/pkg"
	"book/server"
)

func main() {

	db := pkg.DbConnectGorm()
	db.AutoMigrate(&model.Book{})
	server.NewServer(db).StartGin()

}
