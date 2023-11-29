package main

import (
	pb "book/calculator/proto"
	"context"
	"io"
	"log"
	"time"
)

func doMax(c pb.CalculatorServiceClient) {
	log.Println("Starting to do a Client Streaming RPC...")

	stream, err := c.Max(context.Background())
	if err != nil {
		log.Printf("Error while calling Max RPC: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{1, 5, 3, 6, 200, 20}
		for _, number := range numbers {
			log.Printf("Sending number: %v", number)
			err := stream.Send(&pb.MaxRequest{
				Number: number,
			})
			if err != nil {
				log.Printf("Error while sending data to server: %v", err)
				return
			}
			time.Sleep(1 * time.Second)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Printf("Error while closing stream: %v", err)
			return
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break // Exit loop on end of file
			}
			if err != nil {
				log.Printf("Error while receiving data from server: %v", err)
				break // Exit loop on error
			}
			log.Printf("Received a new maimum: %v", res.Result)
		}
		close(waitc)
	}()

	<-waitc
}
