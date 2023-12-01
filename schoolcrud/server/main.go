//package main
//
//import (
//	pb "book/schoolcrud/proto"
//	"context"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"google.golang.org/grpc"
//	"log"
//	"net"
//)
//
//const (
//	mongoURIGrpc       = "mongodb://school:school007@localhost:7700/"
//	databaseNameGrpc   = "grpcschool"
//	collectionNameGrpc = "grpcschool"
//)
//
//var addr string = "0.0.0.0:50051"
//
//type Server struct {
//	pb.SchoolServiceServer
//	collection *mongo.Collection
//}
//
//func main() {
//	clientOptions := options.Client().ApplyURI(mongoURIGrpc)
//	client, err := mongo.Connect(context.Background(), clientOptions)
//	if err != nil {
//		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
//	}
//
//	err = client.Ping(context.Background(), nil)
//	if err != nil {
//		log.Fatalf("Failed to ping MongoDB: %v\n", err)
//	}
//
//	db := client.Database(databaseNameGrpc)
//	collection := db.Collection(collectionNameGrpc)
//
//	lis, err := net.Listen("tcp", addr)
//	if err != nil {
//		log.Fatalf("Failed to listen: %v\n", err)
//	}
//
//	log.Printf("Listening at %s\n", addr)
//
//	s := grpc.NewServer()
//	server := &Server{collection: collection}
//	pb.RegisterSchoolServiceServer(s, server)
//
//	log.Println("|-----------------------------------------------")
//	log.Printf("| gRPC listening on %s:\n", addr)
//	log.Println("|-----------------------------------------------")
//
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("Failed to serve: %v\n", err)
//	}
//}

package servermodel

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "book/schoolcrud/proto"
	"book/schoolcrud/server/handlers"
	"book/schoolcrud/server/repo"
	"book/schoolcrud/server/usecases"
)

const (
	mongoURI       = "mongodb://school:school007@localhost:7700/"
	databaseName   = "grpcschool"
	collectionName = "grpcschool"
	serverAddress  = "0.0.0.0:50051"
)

func main() {
	// MongoDB connection setup
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	// Checking the connection
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// MongoDB collection
	collection := client.Database(databaseName).Collection(collectionName)

	// Repository, UseCases, and Handlers
	schoolRepo := repo.NewSchoolRepository(collection)
	schoolService := usecases.NewSchoolUsecases(schoolRepo)
	schoolHandler := handlers.NewHandlers(schoolService)

	// Start gRPC server
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on %s", serverAddress)

	grpcServer := grpc.NewServer()
	pb.RegisterSchoolServiceServer(grpcServer, schoolHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over %s: %v", serverAddress, err)
	}
}
