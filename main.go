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

	firebaseDB, err := pkg.ConnectFirebaseEmulator()
	if err != nil {
		log.Fatalf("Failed to connect to Firebase: %v", err)
	}

	srv := server.NewServer(gormDB, mongoClient, mongoCollection, firebaseDB)

	if err := srv.StartMongo(context.Background()); err != nil {
		log.Fatalf("Failed to start MongoDB operations: %v", err)
	}

	if err := srv.StartFirebase(); err != nil {
		log.Fatalf("Failed to start Firebase operations: %v", err)
	}

	srv.StartGin()

}
