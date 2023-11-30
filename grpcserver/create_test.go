package grpcserver

import (
	pb "book/schoolcrud/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

func setupInMemoryMongoDB(t *testing.T) (*mongo.Client, *mongo.Collection, func()) {
	mongoServer, err := memongo.Start("4.0.5") // Specify the MongoDB version
	if err != nil {
		t.Fatalf("memongo.Start failed: %s", err)
	}

	opts := options.Client().ApplyURI(mongoServer.URI())
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		mongoServer.Stop()
		t.Fatalf("mongo.Connect failed: %s", err)
	}

	collection := client.Database("grpc").Collection("schools")

	return client, collection, func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Println("Disconnect failed:", err)
		}
		mongoServer.Stop()
	}
}

func TestCreateSchool(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &Server{collection: collection}

	schoolID := primitive.NewObjectID()
	testSchool := &pb.School{
		SchoolId: schoolID.Hex(),
		Name:     "BKK School",
		Address:  "BKK",
		Phone:    "123-456-7890",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := server.CreateSchool(ctx, testSchool)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	returnedID, err := primitive.ObjectIDFromHex(result.Id)
	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, returnedID)
}
