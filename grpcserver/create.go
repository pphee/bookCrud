package grpcserver

import (
	pb "book/schoolcrud/proto"
	"book/schoolcrud/server"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *Server) CreateSchool(ctx context.Context, req *pb.School) (*pb.SchoolId, error) {
	log.Println("CreateSchool")

	if s.collection == nil {
		log.Println("mongoCollectionGrpc is not initialized")
		return nil, status.Errorf(codes.Internal, "mongoCollectionGrpc is not initialized")
	}

	data := main.SchoolItem{
		SchoolId: req.SchoolId,
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
	}

	res, err := s.collection.InsertOne(ctx, data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &pb.SchoolId{
		Id: oid.Hex(),
	}, nil

}
