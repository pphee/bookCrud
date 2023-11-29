package main

import (
	pb "book/greet/proto"
	"context"
	"io"
	"log"
	"time"
)

func doGreetEveryone(c pb.GreetServiceClient) {
	log.Println("Starting to do a Server Streaming RPC...")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC: %v", err)
	}

	reqs := []*pb.GreetRequest{
		{FirstName: "pee"},
		{FirstName: "phee"},
		{FirstName: "james"},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("Sending message: %v", req)
			err := stream.Send(req)
			if err != nil {
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			return
		}

	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error while receiving: %v", err)
				break
			}

			log.Printf("Received: %v", res.Result)
		}
		close(waitc)
	}()

	<-waitc

}
