package main

import (
	pb "book/schoolcrud/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	addr := "0.0.0.0:50051"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	// Create a new gRPC client
	c := pb.NewSchoolServiceClient(conn)

	CreateSchool(c)
	readSchool(c, "6567039a9c6a41d775918e82")
	updateSchool(c, "656702964e4fd2393397c4de")
	deleteSchool(c, "6568640ee28bb0c9bdf82dd5")
	listSchool(c)
}
