package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	pb "book/schoolcrud/proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteSchool(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &Server{collection: collection}

	// Insert a test school document
	oid := primitive.NewObjectID()
	testSchool := bson.M{
		"_id":      oid,
		"schoolId": "1234567890",
		"name":     "Test School",
		"address":  "123 Test Street",
		"phone":    "123-456-7890",
	}
	_, err := collection.InsertOne(context.Background(), testSchool)
	assert.NoError(t, err)

	// Test deletion of the school
	deleteReq := &pb.SchoolId{Id: oid.Hex()}
	_, err = server.DeleteSchool(context.Background(), deleteReq)
	assert.NoError(t, err)

	// Verify deletion
	count, err := collection.CountDocuments(context.Background(), bson.M{"_id": oid})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	nonExistentID := primitive.NewObjectID().Hex()
	_, err = server.DeleteSchool(context.Background(), &pb.SchoolId{Id: nonExistentID})
	assert.Error(t, err)

	_, err = server.DeleteSchool(context.Background(), &pb.SchoolId{Id: "invalid-id"})
	assert.Error(t, err)
}
