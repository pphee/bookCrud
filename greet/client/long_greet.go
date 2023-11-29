package main

import (
	pb "book/greet/proto"
	"context"
	"log"
	"time"
)

func doLongGreet(c pb.GreetServiceClient) {
	log.Println("Starting to do a LongGreet RPC...")

	req := []*pb.GreetRequest{
		{FirstName: "pee"},
		{FirstName: "dans"},
		{FirstName: "james"},
		{FirstName: "bank"},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	for _, req := range req {
		log.Printf("Sending req: %v", req)
		err := stream.Send(req)
		if err != nil {
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)
	}

	log.Printf("LongGreet Response: %v", res)
}
