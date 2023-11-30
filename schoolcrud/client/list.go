package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
)

func listSchool(c pb.SchoolServiceClient) {
	log.Println("--------------------listSchool--------------------")
	stream, err := c.ListSchool(context.Background(), &emptypb.Empty{})

	if err != nil {
		log.Fatalf("Error while calling ListBlogs: %v\n", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Something happened: %v\n", err)
		}

		log.Println(res)
	}
}
