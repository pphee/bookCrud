package client

import (
	pb "book/schoolcrud/proto"
	"context"
	"log"
)

func CreateSchool(c pb.SchoolServiceClient) string {

	log.Println("-------------CreateSchool------------")

	school := &pb.School{
		SchoolId: "1",
		Name:     "School 1",
		Address:  "Address 1",
		Phone:    "Phone 1",
	}

	res, err := c.CreateSchool(context.Background(), school)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	log.Printf("Response from CreateSchool: %v", res.Id)

	return res.Id

}
