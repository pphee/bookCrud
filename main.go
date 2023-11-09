package main

import (
	model "book/bookcrud"
	"book/pkg"
	"book/server"
	"context"
	"log"
)

func main() {

	mongoClient, mongoCollection, err := pkg.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	gormDB := pkg.DbConnectGorm()
	if err := gormDB.AutoMigrate(&model.Book{}); err != nil {
		log.Fatalf("Failed to run auto migration for Book model: %v", err)
	}

	srv := server.NewServer(gormDB, mongoClient, mongoCollection)

	if err := srv.StartMongo(context.Background()); err != nil {
		log.Fatalf("Failed to start MongoDB operations: %v", err)
	}
	srv.StartGin()

}
