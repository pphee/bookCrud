package repomongo

import (
	models "book/studentcrud"
	"book/studentcrud/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestCreateStudent(t *testing.T) {
	mockCollection := &mocks.MockCollection{}
	mockEncryptionService := &mocks.MockEncryptionService{}
	repo := NewStudentRepository(mockCollection, mockEncryptionService)

	student := &models.Student{
		ID:        primitive.NewObjectID(),
		FirstName: "John",
		LastName:  "Doe",
	}

	encryptedFirstName := "encryptedJohn"
	encryptedLastName := "encryptedDoe"

	mockEncryptionService.On("Encrypt", "John").Return(encryptedFirstName, nil)
	mockEncryptionService.On("Encrypt", "Doe").Return(encryptedLastName, nil)

	insertResult := &mongo.InsertOneResult{InsertedID: student.ID}
	mockCollection.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Student"), mock.AnythingOfType("[]*options.InsertOneOptions")).Return(insertResult, nil)

	result, err := repo.Create(context.Background(), student)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, insertResult.InsertedID, result.InsertedID)
	mockEncryptionService.AssertCalled(t, "Encrypt", "John")
	mockEncryptionService.AssertCalled(t, "Encrypt", "Doe")
	assert.Equal(t, encryptedFirstName, student.FirstName)
	assert.Equal(t, encryptedLastName, student.LastName)
}

func TestFindByID(t *testing.T) {
	mockCollection := &mocks.MockCollection{}
	mockEncryptionService := &mocks.MockEncryptionService{}
	repo := NewStudentRepository(mockCollection, mockEncryptionService)

	studentID := primitive.NewObjectID().Hex()
	decryptedFirstName := "John"
	decryptedLastName := "Doe"

	mockEncryptionService.On("Decrypt", mock.Anything).Return(decryptedFirstName, nil).Once()
	mockEncryptionService.On("Decrypt", mock.Anything).Return(decryptedLastName, nil).Once()

	foundStudent := models.Student{
		ID:        primitive.ObjectID{},
		FirstName: "encryptedJohn",
		LastName:  "encryptedDoe",
	}

	singleResult := &mongo.SingleResult{}
	mockCollection.On("FindOne", mock.Anything, mock.AnythingOfType("primitive.M"), mock.Anything).Return(singleResult, nil).Run(func(args mock.Arguments) {
		_ = singleResult.Decode(&foundStudent)
	})

	student, err := repo.FindByID(context.Background(), studentID)

	assert.NoError(t, err)
	assert.Equal(t, decryptedFirstName, student.FirstName)
	assert.Equal(t, decryptedLastName, student.LastName)
}

func TestFindAll(t *testing.T) {
	mockCollection := &mocks.MockCollection{}
	mockEncryptionService := &mocks.MockEncryptionService{}
	repo := NewStudentRepository(mockCollection, mockEncryptionService)

	decryptedFirstName := "John"
	decryptedLastName := "Doe"

	mockEncryptionService.On("Decrypt", mock.Anything).Return(decryptedFirstName, nil).Once()
	mockEncryptionService.On("Decrypt", mock.Anything).Return(decryptedLastName, nil).Once()

	students := []*models.Student{
		{FirstName: "encryptedJohn", LastName: "encryptedDoe"},
		{FirstName: "encryptedJane", LastName: "encryptedDoe"},
	}

	mockCursor := &mocks.MockCursor{Data: students}
	mockCollection.On("Find", mock.Anything, mock.AnythingOfType("primitive.M"), mock.AnythingOfType("[]*options.FindOptions")).Return(mockCursor, nil)

	result, err := repo.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, student := range result {
		assert.Equal(t, decryptedFirstName, student.FirstName)
		assert.Equal(t, decryptedLastName, student.LastName)
	}
}

func TestUpdateStudent(t *testing.T) {
	mockCollection := &mocks.MockCollection{}
	repo := NewStudentRepository(mockCollection, nil)

	studentID := primitive.NewObjectID().Hex()
	updateData := bson.M{"$set": bson.M{"age": 25}}

	updateResult := &mongo.UpdateResult{}
	mockCollection.On("UpdateOne", mock.Anything, mock.AnythingOfType("bson.M"), mock.Anything).Return(updateResult, nil)

	result, err := repo.Update(context.Background(), studentID, updateData)

	assert.NoError(t, err)
	assert.Equal(t, updateResult, result)
}

func TestDeleteStudent(t *testing.T) {
	mockCollection := &mocks.MockCollection{}
	repo := NewStudentRepository(mockCollection, nil)

	studentID := primitive.NewObjectID().Hex()

	deleteResult := &mongo.DeleteResult{}
	mockCollection.On("DeleteOne", mock.Anything, mock.AnythingOfType("bson.M")).Return(deleteResult, nil)

	result, err := repo.Delete(context.Background(), studentID)

	assert.NoError(t, err)
	assert.Equal(t, deleteResult, result)
}
