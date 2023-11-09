package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type IServer interface {
	StartGin()
	StartMongo(ctx context.Context) error
}

type server struct {
	db              *gorm.DB
	app             *gin.Engine
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
}

func NewServer(db *gorm.DB, mongoClient *mongo.Client, mongoCollection *mongo.Collection) IServer {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	return &server{
		app:             app,
		db:              db,
		mongoClient:     mongoClient,
		mongoCollection: mongoCollection,
	}
}

func (s *server) StartGin() {
	api := s.app.Group("/api")
	modules := InitModule(api, s)
	modules.BookModule()
	port := ":8080"
	log.Printf("Server is starting on %v", port)

	// Start the Gin
	if err := s.app.Run(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}

func (s *server) StartMongo(ctx context.Context) error {
	api := s.app.Group("/api")
	modules := InitModule(api, s)
	modules.StudentModule()
	if err := s.mongoClient.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	log.Println("MongoDB has been started successfully")
	return nil
}
