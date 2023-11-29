package server

import (
	pb "book/schoolcrud/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	pb.SchoolServiceServer
	pb.UnimplementedSchoolServiceServer
	mongoClientGrpc     *mongo.Client
	mongoCollectionGrpc *mongo.Collection
}

type SchoolItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	SchoolId string             `bson:"school_id,omitempty"`
	Name     string             `bson:"name"`
	Address  string             `bson:"address"`
	Phone    string             `bson:"phone"`
}

func documentToSchoolItem(data *SchoolItem) *pb.School {
	return &pb.School{
		Id:       data.ID.Hex(),
		SchoolId: data.SchoolId,
		Name:     data.Name,
		Address:  data.Address,
		Phone:    data.Phone,
	}
}
