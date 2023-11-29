package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"log"
)

func deleteSchool(c pb.SchoolServiceClient, id string) {
	log.Println("------------deleteSchool was invoked--------------")
	_, err := c.DeleteSchool(context.Background(), &pb.SchoolId{Id: id})

	if err != nil {
		log.Fatalf("Error happened while deleting: %v\n", err)
	}

	log.Println("School was deleted")
}
