package main

import (
	pb "book/calculator/proto"
	"context"
	"log"
)

func (s *Server) Sum(ctx context.Context, request *pb.SumRequest) (*pb.SumResponse, error) {
	log.Printf("Received: %v", request)
	return &pb.SumResponse{
		Result: request.FirstNumber + request.SecondNumber,
	}, nil
}
