package main

import (
	pb "book/calculator/proto"
	"context"
	"log"
)

func doAvg(c pb.CalculatorServiceClient) {
	log.Println("Starting to do a Server Streaming RPC...")

	stream, err := c.Avg(context.Background())
	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v", err)
	}

	numbers := []int32{1, 2, 3, 4}

	for _, number := range numbers {
		log.Printf("Sending number: %v", number)
		err := stream.Send(&pb.AvgRequest{
			Number: number,
		})
		if err != nil {
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from Calculator RPC: %v", err)
	}

	log.Printf("Response from AVG: %v", res.Result)
}
