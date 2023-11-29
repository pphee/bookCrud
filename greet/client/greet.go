package main

import (
	pb "book/greet/proto"
	"context"
	"log"
)

func doGreet(c pb.GreetServiceClient) {
	log.Println("Starting to do a Unary RPC...")

	res, err := c.Greet(context.Background(), &pb.GreetRequest{
		FirstName: "John",
	})

	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
