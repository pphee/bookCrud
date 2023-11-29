package main

import (
	pb "book/calculator/proto"
	"context"
	"log"
)

func doCalculator(c pb.CalculatorServiceClient) {
	log.Println("Starting to do a Unary RPC...")

	res, err := c.Sum(context.Background(), &pb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 3,
	})

	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
