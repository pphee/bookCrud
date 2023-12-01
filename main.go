package main

import (
	model "book/bookcrud"
	"book/pkg"
	"book/server"
	"context"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoClient, mongoCollection, err := pkg.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	gormDB := pkg.DbConnectGorm()
	if err := gormDB.AutoMigrate(&model.Book{}); err != nil {
		log.Fatalf("Failed to run auto migration for Book model: %v", err)
	}

	firebaseDB, err := pkg.ConnectFirebase()
	if err != nil {
		log.Fatalf("Failed to connect to Firebase: %v", err)
	}

	srv := server.NewServer(gormDB, mongoClient, mongoCollection, firebaseDB)

	if err := srv.StartMongo(context.Background()); err != nil {
		log.Fatalf("Failed to start MongoDB operations: %v", err)
	}

	if err := srv.StartFirebase(context.Background()); err != nil {
		log.Fatalf("Failed to start Firebase operations: %v", err)
	}

	srv.StartGin()

}
