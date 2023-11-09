package handlersmongo

import (
	models "book/studentcrud"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockStudentUsecase struct {
	mock.Mock
}

func (m *mockStudentUsecase) GetAllStudents(ctx context.Context) ([]*models.Student, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Student), args.Error(1)
}

func (m *mockStudentUsecase) CreateStudent(ctx context.Context, student *models.Student) (primitive.ObjectID, error) {
	args := m.Called(ctx, student)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *mockStudentUsecase) DeleteStudent(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockStudentUsecase) GetStudentByID(ctx context.Context, id string) (*models.Student, error) {
	args := m.Called(ctx, id)
	if item := args.Get(0); item != nil {
		return item.(*models.Student), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStudentUsecase) UpdateStudent(ctx context.Context, id string, student models.Student) error {
	args := m.Called(ctx, id, student)
	return args.Error(0)
}

func TestListStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mockStudentUsecase)
	handler := NewStudentHandler(mockUsecase)

	objectId := primitive.NewObjectID()

	expectedStudents := []*models.Student{{ID: objectId, FirstName: "jack", LastName: "q", Age: 20}}

	r := gin.Default()
	r.GET("/students", handler.GetAllStudents)

	tests := []struct {
		description      string
		returnedStudents []*models.Student
		returnedError    error
		expectedCode     int
		expectedBody     string
	}{
		{
			description:      "Get all students",
			returnedStudents: expectedStudents,
			expectedCode:     http.StatusOK,
			expectedBody:     "",
		},
		{
			description:      "Internal server error",
			returnedStudents: nil,
			returnedError:    errors.New("internal server error"),
			expectedCode:     http.StatusInternalServerError,
			expectedBody:     `{"error":"internal server error"}`,
		},
	}
	for _, tc := range tests {
		mockUsecase.ExpectedCalls = nil
		mockUsecase.Calls = nil

		mockUsecase.On("GetAllStudents", mock.Anything).Return(tc.returnedStudents, tc.returnedError)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/students", nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, tc.expectedCode, w.Code)
		if tc.expectedCode == http.StatusOK {
			var response []*models.Student
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.Nil(t, err)
			assert.Equal(t, tc.returnedStudents, response)
		} else {
			assert.Equal(t, tc.expectedBody, w.Body.String())
		}

	}

}

func TestGetStudentByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(mockStudentUsecase)
	handler := NewStudentHandler(mockUseCase)

	objectId := primitive.NewObjectID()

	expectedStudent := &models.Student{ID: objectId, FirstName: "jack", LastName: "q", Age: 20}

	r := gin.Default()
	r.GET("/students/:id", handler.GetStudent)

	tests := []struct {
		description     string
		studentID       string
		mockReturn      *models.Student
		mockError       error
		expectedCode    int
		expectedStudent *models.Student
	}{
		{
			description:     "Successfully retrieved student by ID",
			studentID:       objectId.Hex(),
			mockReturn:      expectedStudent,
			expectedCode:    http.StatusOK,
			expectedStudent: expectedStudent,
		},
		{
			description:  "Student not found",
			studentID:    primitive.NewObjectID().Hex(), // A different ID to simulate not found
			mockError:    errors.New("student not found"),
			expectedCode: http.StatusNotFound,
		},
	}

	// Run the tests
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			mockUseCase.On("GetStudentByID", mock.Anything, tc.studentID).Return(tc.mockReturn, tc.mockError)

			req, _ := http.NewRequest(http.MethodGet, "/students/"+tc.studentID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				var student models.Student
				err := json.Unmarshal(w.Body.Bytes(), &student)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedStudent, &student)
			} else {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, err)
				assert.Equal(t, "Student not found", response["error"])
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestPostStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(mockStudentUsecase)
	handler := NewStudentHandler(mockUseCase)
	validStudent := &models.Student{
		ID:        primitive.NewObjectID(),
		FirstName: "John",
		LastName:  "Doe",
		Age:       20,
	}

	validStudentJSON, _ := json.Marshal(validStudent)
	invalidInputStudentJSON := []byte(`{"FirstName": "", "LastName": "", "Age": 0}`)   // This can still be valid JSON but fails your application's validation
	invalidJSON := []byte(`{"FirstName": "John", "LastName": "Doe", "Age": "twenty"}`) // Invalid JSON, age should be an integer

	r := gin.Default()
	r.POST("/students", handler.PostStudent)

	tests := []struct {
		description  string
		requestBody  []byte
		expectedCode int
		prepare      func()
	}{
		{
			description:  "Successfully created student",
			requestBody:  validStudentJSON,
			expectedCode: http.StatusCreated,
			prepare: func() {
				mockUseCase.On("CreateStudent", mock.Anything, mock.AnythingOfType("*models.Student")).Return(validStudent.ID, nil).Once()
			},
		},
		{
			description:  "Invalid input",
			requestBody:  invalidInputStudentJSON,
			expectedCode: http.StatusBadRequest,
			prepare: func() {
				// No need to set up mock expectations for this case as it should fail validation before calling the usecase
			},
		},
		{
			description:  "Invalid JSON format",
			requestBody:  invalidJSON,
			expectedCode: http.StatusBadRequest,
			prepare: func() {
				// No mock expectations set as the binding will fail before usecase is called
			},
		},
		{
			description:  "Internal server error on creating student",
			requestBody:  validStudentJSON,
			expectedCode: http.StatusInternalServerError,
			prepare: func() {
				mockUseCase.On("CreateStudent", mock.Anything, mock.AnythingOfType("*models.Student")).Return(primitive.NilObjectID, errors.New("internal server error")).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			mockUseCase.ExpectedCalls = nil
			mockUseCase.Calls = nil

			tc.prepare()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusCreated {
				var studentresp models.Student
				err := json.Unmarshal(w.Body.Bytes(), &studentresp)
				assert.NoError(t, err)
				assert.Equal(t, validStudent.ID.Hex(), studentresp.ID.Hex())
			}
		})
	}
}

func TestDeleteStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(mockStudentUsecase)
	handler := NewStudentHandler(mockUseCase)

	student := &models.Student{
		FirstName: "pee",
		LastName:  "dans",
		Age:       20,
	}

	r := gin.Default()
	r.DELETE("/students/:id", handler.DeleteStudent)

	tests := []struct {
		description  string
		student      *models.Student
		expectedCode int
		prepare      func()
	}{
		{
			description:  "Successfully deleted student",
			student:      student,
			expectedCode: http.StatusOK,
			prepare: func() {
				mockUseCase.On("DeleteStudent", mock.Anything).Return(nil)
			},
		},
		{
			description:  "Student not found or unable to delete",
			student:      student,
			expectedCode: http.StatusNotFound,
			prepare: func() {
				mockUseCase.On("DeleteStudent", mock.Anything).Return(errors.New("student not found"))
			},
		},
	}

	for _, tc := range tests {
		mockUseCase.ExpectedCalls = nil
		mockUseCase.Calls = nil

		// Act
		tc.prepare()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/students/1", nil)
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, tc.expectedCode, w.Code)
	}
}

func TestUpdateStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(mockStudentUsecase)
	handler := NewStudentHandler(mockUseCase)

	studentID := primitive.NewObjectID().Hex()

	validStudent := models.Student{
		FirstName: "John",
		LastName:  "Doe",
		Age:       20,
	}

	r := gin.Default()
	r.PUT("/students/:id", handler.UpdateStudent)

	tests := []struct {
		description  string
		studentID    string
		student      models.Student
		expectedCode int
		prepare      func(studentID string, student models.Student)
		want         string
	}{
		{
			description:  "Successfully updated student",
			studentID:    studentID,
			student:      validStudent,
			expectedCode: http.StatusOK,
			prepare: func(id string, s models.Student) {
				mockUseCase.On("UpdateStudent", mock.Anything, id, s).Return(nil)
			},
			want: `{"message":"Student updated successfully"}`,
		},
		{
			description:  "Internal server error on updating student",
			studentID:    studentID,
			student:      validStudent,
			expectedCode: http.StatusInternalServerError,
			prepare: func(id string, s models.Student) {
				mockUseCase.On("UpdateStudent", mock.Anything, id, s).Return(errors.New("internal server error"))
			},
			want: "internal server error",
		},
		{
			description: "Missing required fields",
			studentID:   studentID,
			student: models.Student{
				FirstName: "",
				LastName:  "",
				Age:       0,
			},
			expectedCode: http.StatusBadRequest,
			prepare: func(studentID string, student models.Student) {
				// Again, no need for mock setup
			},
			want: "invalid input, missing fields",
		},
		{
			description:  "Invalid JSON input",
			studentID:    studentID,
			student:      models.Student{}, // an empty student object, which will fail JSON binding
			expectedCode: http.StatusBadRequest,
			prepare: func(studentID string, student models.Student) {
				// No need to set up mock expectations since the request should fail before it hits the use case
			},
			want: "invalid input, missing fields",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			mockUseCase.ExpectedCalls = nil
			mockUseCase.Calls = nil

			tc.prepare(tc.studentID, tc.student)

			body, _ := json.Marshal(tc.student)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/students/%s", tc.studentID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if w.Code == http.StatusOK {
				assert.JSONEq(t, tc.want, w.Body.String())
			} else {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.want, response["error"])
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}
