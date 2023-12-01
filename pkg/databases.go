package pkg

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func DbConnectGorm() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database.")
	}
	return db
}

const (
	mongoURI       = "mongodb://cloud:cloud007@localhost:27017/"
	databaseName   = "cloudbook"
	collectionName = "cloudbook"
)

func ConnectMongoDB() (*mongo.Client, *mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	return client, collection, nil
}

func ConnectFirebase() (*firestore.Client, error) {
	firestoreEmulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if firestoreEmulatorHost == "" {
		return nil, fmt.Errorf("FIRESTORE_EMULATOR_HOST environment variable is not set")
	}

	project := os.Getenv("PROJECT")
	if project == "" {
		return nil, fmt.Errorf("PROJECT environment variable is not set")
	}

	client, err := firestore.NewClient(context.Background(), project)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return client, nil
}
