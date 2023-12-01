package repo

import (
	pb "book/schoolcrud/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
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

type SchoolItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	SchoolId string             `bson:"school_id"`
	Name     string             `bson:"name"`
	Address  string             `bson:"address"`
	Phone    string             `bson:"phone"`
}

func TestCreateSchool(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &mongoSchoolRepository{collection: collection}

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
}

func TestListSchool(t *testing.T) {
	// Setup in-memory MongoDB
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &mongoSchoolRepository{collection: collection}

	// Insert mock data into the collection
	for i := 0; i < 3; i++ {
		school := &pb.School{
			SchoolId: primitive.NewObjectID().Hex(),
			Name:     "Test School " + strconv.Itoa(i),
			Address:  "123 Test Street",
			Phone:    "123-456-7890",
		}
		_, err := collection.InsertOne(context.Background(), school)
		assert.NoError(t, err)
	}

	schools, err := server.ListSchools(context.Background())
	assert.NoError(t, err)

	assert.Len(t, schools, 3)
}

func TestGetSchool(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &mongoSchoolRepository{collection: collection}

	testSchool := &SchoolItem{
		ID:       primitive.NewObjectID(),
		SchoolId: "1",
		Name:     "BKK School",
		Address:  "BKK",
		Phone:    "123-456-7890",
	}
	_, err := collection.InsertOne(context.Background(), testSchool)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := server.GetSchool(ctx, testSchool.ID.Hex())

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, testSchool.Name, result.Name)
		assert.Equal(t, testSchool.Address, result.Address)
		assert.Equal(t, testSchool.Phone, result.Phone)
	}
}

func TestUpdateSchool(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &mongoSchoolRepository{collection: collection}

	oid := primitive.NewObjectID()
	testSchool := &SchoolItem{
		ID:       oid,
		SchoolId: "1",
		Name:     "Yala School",
		Address:  "Yala",
		Phone:    "073212761",
	}
	_, err := collection.InsertOne(context.Background(), testSchool)
	assert.NoError(t, err)

	updateReq := &pb.School{
		Id:       oid.Hex(),
		SchoolId: "2",
		Name:     "Thailand School",
		Address:  "Thailand",
		Phone:    "073203778",
	}

	err = server.UpdateSchool(context.Background(), updateReq)
	assert.NoError(t, err)

	var getUpdatedReqSchool SchoolItem
	err = collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&getUpdatedReqSchool)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Name, getUpdatedReqSchool.Name)
	assert.Equal(t, updateReq.Address, getUpdatedReqSchool.Address)
	assert.Equal(t, updateReq.Phone, getUpdatedReqSchool.Phone)
}

func TestDeleteSchool(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	server := &mongoSchoolRepository{collection: collection}

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
	err = server.DeleteSchool(context.Background(), deleteReq.Id)
	assert.NoError(t, err)

	// Verify deletion
	count, err := collection.CountDocuments(context.Background(), bson.M{"_id": oid})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	nonExistentID := primitive.NewObjectID().Hex()
	err = server.DeleteSchool(context.Background(), nonExistentID) // Pass only the ID
	assert.Error(t, err)

	err = server.DeleteSchool(context.Background(), "invalid-id") // Pass only the ID
	assert.Error(t, err)
}
