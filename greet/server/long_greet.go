package main

import (
	pb "book/greet/proto"
	"fmt"
	"io"
	"log"
)

func (s *Server) LongGreet(stream pb.GreetService_LongGreetServer) error {

	log.Println("LongGreet Function was invoked")

	res := ""

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			// We have finished reading the client stream.
			return stream.SendAndClose(&pb.GreetResponse{
				Result: res,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		log.Printf("Received: %v", req)
		res += fmt.Sprintf("Hello %s!\n", req.FirstName)
	}
	return nil
}
