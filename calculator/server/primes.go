package main

import (
	pb "book/calculator/proto"
	"log"
)

func (s *Server) Prime(request *pb.PrimeRequest, stream pb.CalculatorService_PrimeServer) error {
	log.Printf("Received prime : %v", request)

	number := request.Number
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			err := stream.Send(&pb.PrimeResponse{
				Result: divisor,
			})

			number /= divisor
			if err != nil {
				return err
			}
		} else {
			divisor++
		}
	}

	return nil
}
