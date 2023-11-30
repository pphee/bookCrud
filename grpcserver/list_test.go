package grpcserver

import (
	"book/schoolcrud/server"
	"context"
	"strconv"
	"testing"

	pb "book/schoolcrud/proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type mockListSchoolServer struct {
	pb.SchoolService_ListSchoolServer
	sendMsg []*pb.School
}

func (m *mockListSchoolServer) Send(school *pb.School) error {
	m.sendMsg = append(m.sendMsg, school)
	return nil
}

func TestListSchool(t *testing.T) {
	// Setup in-memory MongoDB
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &Server{collection: collection}

	// Insert mock data into the collection
	for i := 0; i < 3; i++ {
		school := &main.SchoolItem{
			SchoolId: primitive.NewObjectID().Hex(),
			Name:     "Test School " + strconv.Itoa(i),
			Address:  "123 Test Street",
			Phone:    "123-456-7890",
		}
		_, err := collection.InsertOne(context.Background(), school)
		assert.NoError(t, err)
	}

	mockStream := &mockListSchoolServer{sendMsg: make([]*pb.School, 0)}

	err := server.ListSchool(&emptypb.Empty{}, mockStream)
	assert.NoError(t, err)

	assert.Len(t, mockStream.sendMsg, 3)
}
