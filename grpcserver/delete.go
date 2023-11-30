package grpcserver

import (
	pb "book/schoolcrud/proto"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *Server) DeleteSchool(ctx context.Context, req *pb.SchoolId) (*emptypb.Empty, error) {
	log.Printf("DeleteSchool was invoked with %v\n", req)

	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"School was not found",
		)
	}

	return &emptypb.Empty{}, nil
}
