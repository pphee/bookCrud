package main

import (
	pb "book/calculator/proto"
	"io"
	"log"
)

func (s *Server) Avg(stream pb.CalculatorService_AvgServer) error {
	log.Println("Received Avg")

	var sum int32 = 0
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AvgResponse{
				Result: float64(sum) / float64(count),
			})
		}
		if err != nil {
			return err
		}

		log.Printf("Received number : %v", req.Number)
		sum += req.Number
		count++
	}
}
