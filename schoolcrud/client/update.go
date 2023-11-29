package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"log"
)

func updateSchool(c pb.SchoolServiceClient, id string) {
	log.Println("-------------updateSchool------------")
	newBlog := &pb.School{
		Id:       id,
		SchoolId: "updated-school-id",
		Name:     "JQK",
		Address:  "anime",
		Phone:    "1999",
	}

	_, err := c.UpdateSchool(context.Background(), newBlog)

	if err != nil {
		log.Printf("Error happened while updating: %v\n", err)
	}

	log.Println("School was updated")
}
