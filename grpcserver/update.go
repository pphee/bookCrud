package grpcserver

//
//import (
//	pb "book/schoolcrud/proto"
//	"book/schoolcrud/server"
//	"context"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"google.golang.org/protobuf/types/known/emptypb"
//	"log"
//)
//
//func (s *Server) UpdateSchool(ctx context.Context, req *pb.School) (*emptypb.Empty, error) {
//	log.Println("UpdateSchool", req)
//
//	oid, err := primitive.ObjectIDFromHex(req.Id)
//	if err != nil {
//		return nil, status.Errorf(
//			codes.InvalidArgument,
//			"Cannot parse ID",
//		)
//	}
//
//	data := &main.SchoolItem{
//		SchoolId: req.SchoolId,
//		Name:     req.Name,
//		Address:  req.Address,
//		Phone:    req.Phone,
//	}
//	res, err := s.collection.UpdateOne(
//		ctx,
//		bson.M{"_id": oid},
//		bson.M{"$set": data},
//	)
//
//	if err != nil {
//		return nil, status.Errorf(
//			codes.Internal,
//			"Could not update",
//		)
//	}
//
//	if res.MatchedCount == 0 {
//		return nil, status.Errorf(
//			codes.NotFound,
//			"Cannot find blog with ID",
//		)
//	}
//
//	return &emptypb.Empty{}, nil
//}
