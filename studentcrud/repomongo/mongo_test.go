package repomongo_test

import (
	models "book/studentcrud"
	"book/studentcrud/repomongo"
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	collection := client.Database("test").Collection("students")

	return client, collection, func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Println("Disconnect failed:", err)
		}
		mongoServer.Stop()
	}
}

func TestCreateStudent(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	student := models.Student{
		ID:        primitive.NewObjectID(),
		FirstName: "John",
		LastName:  "Doe",
		Age:       21,
	}

	repo := repomongo.NewStudentRepository(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, insertErr := repo.Create(ctx, &student)
	if insertErr != nil {
		t.Fatalf("Failed to create student: %v", insertErr)
	}
	assert.NotNil(t, result)
	assert.Equal(t, student.ID, result.InsertedID.(primitive.ObjectID), "The inserted ID should match the student ID")
}

func TestFindByID(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	repo := repomongo.NewStudentRepository(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	studentID := "6549e760c7623d4e258a8397"
	objID, err := primitive.ObjectIDFromHex(studentID)
	assert.NoError(t, err)

	_, err = collection.InsertOne(ctx, bson.M{"_id": objID})
	assert.NoError(t, err)

	student, err := repo.FindByID(ctx, studentID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, objID, student.ID)
}

func TestFindAll(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	studentsToInsert := []interface{}{
		models.Student{
			ID:        primitive.NewObjectID(),
			FirstName: "John",
			LastName:  "Doe",
		},
		models.Student{
			ID:        primitive.NewObjectID(),
			FirstName: "Jane",
			LastName:  "Smith",
		},
	}

	_, err := collection.InsertMany(context.Background(), studentsToInsert)
	assert.NoError(t, err)

	repo := repomongo.NewStudentRepository(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	students, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, students)
	assert.True(t, len(students) > 0, "there should be students in the collection")
	assert.Equal(t, studentsToInsert[0].(models.Student).ID.Hex(), students[0].ID.Hex())
	assert.Equal(t, studentsToInsert[0].(models.Student).FirstName, students[0].FirstName)

}

func TestUpdate(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	repo := repomongo.NewStudentRepository(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex("507f191e810c19729de860ea")
	assert.NoError(t, err)

	_, err = collection.InsertOne(ctx, bson.M{"_id": objID, "FirstName": "John"})
	assert.NoError(t, err)

	studentID := "507f191e810c19729de860ea"
	updateData := bson.M{
		"FirstName": "Jane",
	}

	result, err := repo.Update(ctx, studentID, updateData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ModifiedCount)
}

func TestDelete(t *testing.T) {
	_, collection, cleanup := setupInMemoryMongoDB(t)
	defer cleanup()

	repo := repomongo.NewStudentRepository(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex("507f191e810c19729de860ea")
	assert.NoError(t, err)

	_, err = collection.InsertOne(ctx, bson.M{"_id": objID, "FirstName": "John"})
	assert.NoError(t, err)

	studentID := "507f191e810c19729de860ea"
	result, err := repo.Delete(ctx, studentID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.DeletedCount)
}
