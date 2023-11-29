package main

import (
	pb "book/calculator/proto"
	"io"
	"log"
)

func (s *Server) Max(stream pb.CalculatorService_MaxServer) error {
	log.Println("Received Max")

	var max int32 = 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			// End of client stream
			return nil
		}
		if err != nil {
			log.Printf("Error while reading client stream: %v", err)
			return err // Return error instead of exiting the entire server
		}

		if req.Number > max {
			max = req.Number
		}

		err = stream.Send(&pb.MaxResponse{
			Result: max,
		})

		if err != nil {
			log.Printf("Error while sending data to client: %v", err)
			return err // Return error instead of exiting the entire server
		}

		log.Printf("Received: %v", req)
	}
}
