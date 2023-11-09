package usecasesmongo_test

import (
	models "book/studentcrud"
	"book/studentcrud/usecasesmongo"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, student *models.Student) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, student)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *mockRepository) FindByID(ctx context.Context, id string) (models.Student, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Student), args.Error(1)
}

func (m *mockRepository) FindAll(ctx context.Context) ([]*models.Student, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Student), args.Error(1)
}

func (m *mockRepository) Update(ctx context.Context, id string, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, id, update)
	result, _ := args.Get(0), args.Error(1)
	// You need to check if args.Get(0) is actually a *mongo.UpdateResult before casting
	updateResult, ok := result.(*mongo.UpdateResult)
	if !ok {
		return nil, args.Error(1)
	}
	return updateResult, args.Error(1)
}

func (m *mockRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, id)
	result, _ := args.Get(0), args.Error(1)
	deleteResult, ok := result.(*mongo.DeleteResult)
	if !ok {
		return nil, args.Error(1)
	}
	return deleteResult, args.Error(1)
}

func TestCreateStudent(t *testing.T) {
	mockRepo := new(mockRepository)
	useCase := usecasesmongo.NewStudentUseCase(mockRepo)
	ctx := context.TODO()

	newStudent := models.Student{
		ID:        primitive.NewObjectID(),
		FirstName: "John",
		LastName:  "Doe",
		Age:       20,
	}
	insertOneResult := &mongo.InsertOneResult{
		InsertedID: newStudent.ID,
	}

	mockRepo.On("Create", ctx, newStudent).Return(insertOneResult, nil)

	studentID, err := useCase.CreateStudent(ctx, &newStudent)

	assert.NoError(t, err)
	assert.Equal(t, newStudent.ID, studentID)
	mockRepo.AssertExpectations(t)
}

func TestGetStudentByID(t *testing.T) {
	mockRepo := new(mockRepository)
	useCase := usecasesmongo.NewStudentUseCase(mockRepo)
	ctx := context.TODO()

	studentID := primitive.NewObjectID().Hex()
	expectedStudent := models.Student{
		ID:        primitive.ObjectID{},
		FirstName: "John",
		LastName:  "Doe",
		Age:       20,
	}

	mockRepo.On("FindByID", ctx, studentID).Return(expectedStudent, nil)

	student, err := useCase.GetStudentByID(ctx, studentID)

	assert.NoError(t, err)
	assert.Equal(t, expectedStudent, student)
	mockRepo.AssertExpectations(t)
}

func TestDeleteStudent(t *testing.T) {
	mockRepo := new(mockRepository)
	useCase := usecasesmongo.NewStudentUseCase(mockRepo)
	ctx := context.TODO()

	studentID := primitive.NewObjectID().Hex()
	deleteResult := &mongo.DeleteResult{
		DeletedCount: 1,
	}

	mockRepo.On("Delete", ctx, studentID).Return(deleteResult, nil)

	err := useCase.DeleteStudent(ctx, studentID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStudent(t *testing.T) {
	mockRepo := new(mockRepository)
	useCase := usecasesmongo.NewStudentUseCase(mockRepo)
	ctx := context.TODO()

	studentID := primitive.NewObjectID().Hex()
	updatedStudent := models.Student{
		FirstName: "UpdatedFirstName",
		LastName:  "UpdatedLastName",
		Age:       21,
	}

	updateResult := &mongo.UpdateResult{}

	mockRepo.On("Update", ctx, studentID, updatedStudent).Return(updateResult, nil)

	err := useCase.UpdateStudent(ctx, studentID, updatedStudent)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAllStudents(t *testing.T) {
	mockRepo := new(mockRepository)
	useCase := usecasesmongo.NewStudentUseCase(mockRepo)
	ctx := context.TODO()

	expectedStudents := []models.Student{
		{
			ID:        primitive.NewObjectID(),
			FirstName: "phee",
			LastName:  "dans",
			Age:       20,
		},
		{
			ID:        primitive.NewObjectID(),
			FirstName: "pruck",
			LastName:  "dans",
			Age:       21,
		},
	}

	mockRepo.On("FindAll", ctx).Return(expectedStudents, nil)
	student, err := useCase.GetAllStudents(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedStudents, student)
	mockRepo.AssertExpectations(t)
}
