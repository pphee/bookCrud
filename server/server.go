package server

import (
	"book/bookcrud/repo"
	"book/generator"
	pb "book/schoolcrud/proto"
	"book/studentcrud/repomongo"
	"cloud.google.com/go/firestore"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type IServer interface {
	StartGin()
	StartMongo(ctx context.Context) error
	StartFirebase(ctx context.Context) error
	StartGrpc(ctx context.Context) error
}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type Dog struct {
	Name    string                 `firestore:"name"`
	Age     int                    `firestore:"age"`
	Human   *firestore.DocumentRef `firestore:"human"`
	Created time.Time              `firestore:"created"`
}

func NewDog(name string, age int, human *firestore.DocumentRef) Dog {
	return Dog{
		Name:    name,
		Age:     age,
		Human:   human,
		Created: time.Now(),
	}
}

func (d *Dog) ID() string {
	return hash(d.Name)
}

type Human struct {
	Name string `firestore:"name"`
}

func (h *Human) ID() string {
	return hash(h.Name)
}

type server struct {
	db                  *gorm.DB
	app                 *gin.Engine
	mongoClient         *mongo.Client
	mongoCollection     *mongo.Collection
	studentRepo         repomongo.IStudentRepository
	bookRepo            repo.IBookRepository
	firebaseClient      *firestore.Client
	mongoClientGrpc     *mongo.Client
	mongoCollectionGrpc *mongo.Collection
}

func NewServer(db *gorm.DB, mongoClient *mongo.Client, mongoCollection *mongo.Collection, firebaseClient *firestore.Client, mongoClientGrpc *mongo.Client, mongoCollectionGrpc *mongo.Collection) IServer {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	encryptionKey := []byte("jack-queen-kingpok92deng") // Use your actual encryption key
	studentRepo := repomongo.NewStudentRepository(mongoCollection, encryptionKey)

	bookRepo := repo.NewGormStore(db)

	return &server{
		app:                 app,
		db:                  db,
		mongoClient:         mongoClient,
		mongoCollection:     mongoCollection,
		studentRepo:         studentRepo,
		bookRepo:            bookRepo,
		firebaseClient:      firebaseClient,
		mongoClientGrpc:     mongoClientGrpc,
		mongoCollectionGrpc: mongoCollectionGrpc,
	}
}

func (s *server) StartGin() {
	api := s.app.Group("/api")
	modules := InitModule(api, s, s.firebaseClient)
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
	modules := InitModule(api, s, s.firebaseClient)
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

func (s *server) StartFirebase(ctx context.Context) error {
	api := s.app.Group("/api")
	modules := InitModule(api, s, s.firebaseClient)
	modules.TeacherModule()
	me := Human{Name: "me"}
	meDocRef := s.firebaseClient.Collection("humans").Doc(me.ID())
	if _, err := meDocRef.Set(ctx, me); err != nil {
		log.Fatalf("failed to add human to Firestore: %v", err)
	}

	freddie := NewDog("Freddie", 2, meDocRef)
	freddieDocRef := s.firebaseClient.Collection("dogs").Doc(freddie.ID())
	if _, err := freddieDocRef.Set(ctx, freddie); err != nil {
		log.Fatalf("failed to add dog to Firestore: %v", err)
	}

	log.Println("Firestore has been started successfully")
	return nil
}

type Server struct {
	pb.SchoolServiceServer
	mongoClientGrpc     *mongo.Client
	mongoCollectionGrpc *mongo.Collection
}

func (s *server) StartGrpc(ctx context.Context) error {
	g := grpc.NewServer()
	pb.RegisterSchoolServiceServer(g, &Server{
		mongoClientGrpc:     s.mongoClientGrpc,
		mongoCollectionGrpc: s.mongoCollectionGrpc,
	})
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		if err := g.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Println("gRPC client connected successfully")
	return nil
}
