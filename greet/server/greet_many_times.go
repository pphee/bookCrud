package main

import (
	pb "book/greet/proto"
	"fmt"
	"log"
)

func (s *Server) GreetManyTimes(request *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function was invoked with %v\n", request)

	for i := 0; i < 10; i++ {
		result := fmt.Sprintf("Fuck %s number %d", request.FirstName, i)
		response := &pb.GreetResponse{
			Result: result,
		}

		err := stream.Send(response)
		if err != nil {
			return err
		}
	}

	return nil

}
