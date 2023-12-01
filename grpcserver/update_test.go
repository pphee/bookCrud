package grpcserver

//
//import (
//	"book/schoolcrud/server"
//	"context"
//	"testing"
//
//	pb "book/schoolcrud/proto"
//	"github.com/stretchr/testify/assert"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//func TestUpdateSchool(t *testing.T) {
//	_, collection, cleanup := setupInMemoryMongoDB(t)
//	defer cleanup()
//
//	server := &Server{collection: collection}
//
//	oid := primitive.NewObjectID()
//	testSchool := &main.SchoolItem{
//		ID:       oid,
//		SchoolId: oid.Hex(),
//		Name:     "Yala School",
//		Address:  "Yala",
//		Phone:    "073212761",
//	}
//	_, err := collection.InsertOne(context.Background(), testSchool)
//	assert.NoError(t, err)
//
//	updateReq := &pb.School{
//		Id:      oid.Hex(),
//		Name:    "Thailand School",
//		Address: "Thailand",
//		Phone:   "073203778",
//	}
//
//	_, err = server.UpdateSchool(context.Background(), updateReq)
//	assert.NoError(t, err)
//
//	var getUpdatedReqSchool main.SchoolItem
//	err = collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&getUpdatedReqSchool)
//	assert.NoError(t, err)
//	assert.Equal(t, updateReq.Name, getUpdatedReqSchool.Name)
//	assert.Equal(t, updateReq.Address, getUpdatedReqSchool.Address)
//	assert.Equal(t, updateReq.Phone, getUpdatedReqSchool.Phone)
//}
