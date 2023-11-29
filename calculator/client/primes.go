package main

import (
	pb "book/calculator/proto"
	"context"
	"io"
	"log"
)

func doPrimes(c pb.CalculatorServiceClient) {
	log.Println("Starting to do a Server Streaming RPC...")

	req := &pb.PrimeRequest{
		Number: 120,
	}

	resStream, err := c.Prime(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Prime RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from Prime: %v", msg.Result)
	}

}
