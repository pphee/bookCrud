package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *Server) GetSchool(ctx context.Context, req *pb.SchoolId) (*pb.School, error) {
	log.Println("GetSchool")

	oid, err := primitive.ObjectIDFromHex(req.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	data := &SchoolItem{}
	filter := bson.M{"_id": oid}

	res := s.collection.FindOne(ctx, filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find school with specified ID: %v", err),
		)
	}

	return documentToSchoolItem(data), nil
}
