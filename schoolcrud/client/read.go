package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"log"
	"time"
)

func readSchool(c pb.SchoolServiceClient, id string) {
	log.Println("-------------ReadSchool------------")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SchoolId{
		Id: id,
	}
	res, err := c.GetSchool(ctx, req)
	if err != nil {
		log.Printf("Unexpected error: %v", err)
		return
	}
	log.Printf("Response from ReadSchool: %v", res)
}
