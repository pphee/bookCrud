package server

import (
	"book/bookcrud/repo"
	"book/generator"
	"book/studentcrud/repomongo"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
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
	StartFirebase() error
}

type server struct {
	db              *gorm.DB
	app             *gin.Engine
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
	studentRepo     repomongo.IStudentRepository
	bookRepo        repo.IBookRepository
	firebaseApp     *firebase.App
}

func NewServer(db *gorm.DB, mongoClient *mongo.Client, mongoCollection *mongo.Collection, firebaseApp *firebase.App) IServer {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	encryptionKey := []byte("jack-queen-kingpok92deng") // Use your actual encryption key
	studentRepo := repomongo.NewStudentRepository(mongoCollection, encryptionKey)

	bookRepo := repo.NewGormStore(db)

	return &server{
		app:             app,
		db:              db,
		mongoClient:     mongoClient,
		mongoCollection: mongoCollection,
		studentRepo:     studentRepo,
		bookRepo:        bookRepo,
		firebaseApp:     firebaseApp,
	}
}

func (s *server) StartGin() {
	api := s.app.Group("/api")
	modules := InitModule(api, s)
	modules.BookModule()
	port := ":8080"
	log.Printf("Server is starting on %v", port)

	//book, err := generator.GenerateBooks()
	//if err != nil {
	//	log.Fatalf("failed to generate books: %v", err)
	//}
	//for _, v := range book {
	//	fmt.Println(v)
	//	_, err = s.bookRepo.New(&v)
	//	if err != nil {
	//		log.Printf("failed to insert book: %v", err)
	//		continue
	//	}
	//}

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

	student, err := generator.GenerateStudents()
	if err != nil {
		return fmt.Errorf("failed to generate users: %v", err)
	}

	for _, v := range student {
		fmt.Println(student)
		_, err = s.studentRepo.Create(ctx, &v)
		if err != nil {
			return fmt.Errorf("failed to insert student: %v", err)
		}
	}

	studentID := student[0].ID.Hex()
	_, err = s.studentRepo.FindByID(ctx, studentID)

	if err != nil {
		return fmt.Errorf("failed to find student by ID: %v", err)
	}

	if err != nil {
		return fmt.Errorf("failed to insert student: %v", err)
	}
	log.Println("MongoDB has been started successfully")
	return nil
}

func (s *server) StartFirebase() error {
	ctx := context.Background()
	client, err := s.firebaseApp.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to Firestore: %w", err)
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("failed to close Firestore client: %v", err)
		}
	}(client)

	log.Println("firebase has been started successfully")

	return nil
}
