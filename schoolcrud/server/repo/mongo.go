package repo

import (
	pb "book/schoolcrud/proto"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type ISchoolRepository interface {
	CreateSchool(ctx context.Context, school *pb.School) (*pb.SchoolId, error)
	ListSchools(ctx context.Context) ([]*pb.School, error)
	GetSchool(ctx context.Context, id string) (*pb.School, error)
	UpdateSchool(ctx context.Context, school *pb.School) error
	DeleteSchool(ctx context.Context, id string) error
}

type mongoSchoolRepository struct {
	collection *mongo.Collection
}

func NewSchoolRepository(collection *mongo.Collection) ISchoolRepository {
	return &mongoSchoolRepository{
		collection: collection,
	}
}

func (m *mongoSchoolRepository) CreateSchool(ctx context.Context, school *pb.School) (*pb.SchoolId, error) {
	res, err := m.collection.InsertOne(ctx, school)
	if err != nil {
		log.Printf("Failed to insert school: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Cannot convert to OID")
	}

	return &pb.SchoolId{Id: oid.Hex()}, nil
}

func (m *mongoSchoolRepository) ListSchools(ctx context.Context) ([]*pb.School, error) {
	var schools []*pb.School

	cur, err := m.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unknown internal error: %v", err)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	for cur.Next(ctx) {
		var school pb.School
		err := cur.Decode(&school)
		if err != nil {
			log.Printf("Error while decoding data from MongoDB: %v", err)
			continue
		}
		schools = append(schools, &school)
	}

	if err := cur.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Unknown internal error: %v", err)

	}

	return schools, nil
}

func (m *mongoSchoolRepository) GetSchool(ctx context.Context, id string) (*pb.School, error) {
	var school pb.School

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse ID")
	}

	filter := bson.M{"_id": oid}

	err = m.collection.FindOne(ctx, filter).Decode(&school)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Cannot find school with specified ID: %v", err)
	}

	return &school, nil
}

func (m *mongoSchoolRepository) UpdateSchool(ctx context.Context, school *pb.School) error {
	oid, err := primitive.ObjectIDFromHex(school.Id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Cannot parse ID: %v", err)
	}

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": school},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "Could not update: %v", err)
	}

	if res.MatchedCount == 0 {
		return status.Errorf(codes.NotFound, "Cannot find school with ID")
	}

	return nil
}

func (m *mongoSchoolRepository) DeleteSchool(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Cannot parse ID: %v", err)
	}

	res, err := m.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return status.Errorf(codes.Internal, "Could not delete: %v", err)
	}

	if res.DeletedCount == 0 {
		return status.Errorf(codes.NotFound, "Cannot find school with ID")
	}

	return nil
}
