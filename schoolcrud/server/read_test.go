package main

import (
	pb "book/schoolcrud/proto"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestGetSchool(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &Server{collection: collection}

	testSchoolItem := &SchoolItem{
		ID:       primitive.NewObjectID(),
		SchoolId: "1",
		Name:     "Test School",
		Address:  "123 Test Street",
		Phone:    "123-456-7890",
	}
	_, err := collection.InsertOne(context.Background(), testSchoolItem)
	fmt.Println("TestGetSchool", testSchoolItem)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := server.GetSchool(ctx, &pb.SchoolId{Id: testSchoolItem.ID.Hex()})

	assert.NoError(t, err)
	if assert.NotNil(t, result) { // Check result is not nil before accessing its fields
		assert.Equal(t, testSchoolItem.Name, result.Name)
		assert.Equal(t, testSchoolItem.Address, result.Address)
		assert.Equal(t, testSchoolItem.Phone, result.Phone)
	}
}
