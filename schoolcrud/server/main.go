package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	mongoURIGrpc       = "mongodb://school:school007@localhost:7700/"
	databaseNameGrpc   = "grpcschool"
	collectionNameGrpc = "grpcschool"
)

var addr string = "0.0.0.0:50051"

// Server struct to hold the collection
type Server struct {
	collection *mongo.Collection
	pb.SchoolServiceServer
}

func main() {
	clientOptions := options.Client().ApplyURI(mongoURIGrpc)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v\n", err)
	}

	db := client.Database(databaseNameGrpc)
	collection := db.Collection(collectionNameGrpc)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	log.Printf("Listening at %s\n", addr)

	s := grpc.NewServer()
	server := &Server{collection: collection}
	pb.RegisterSchoolServiceServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
