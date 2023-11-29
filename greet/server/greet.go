package main

import (
	pb "book/greet/proto"
	"context"
	"log"
)

func (s *Server) Greet(ctx context.Context, request *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Received: %v", request)
	return &pb.GreetResponse{
		Result: "Hello " + request.FirstName,
	}, nil
}
