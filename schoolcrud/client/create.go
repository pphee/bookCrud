package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"log"
	"time"
)

func CreateSchool(c pb.SchoolServiceClient) string {
	log.Println("-------------CreateSchool------------")

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	school := &pb.School{
		SchoolId: "1",
		Name:     "School 1",
		Address:  "Address 1",
		Phone:    "Phone 1",
	}

	res, err := c.CreateSchool(ctx, school)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	log.Printf("Response from CreateSchool: %v", res.Id)

	return res.Id
}
