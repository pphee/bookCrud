package handlersfirebase_test

import (
	models "book/teachercrud"
	"book/teachercrud/handlersfirebase"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTeacherUseCases struct {
	Teachers map[string]*models.Teacher
}

func NewMockTeacherUseCases() *MockTeacherUseCases {
	return &MockTeacherUseCases{
		Teachers: make(map[string]*models.Teacher),
	}
}

func (m *MockTeacherUseCases) AddTeacher(ctx context.Context, teacher *models.Teacher) error {
	teacher.HashedID = teacher.HashID()
	m.Teachers[teacher.HashedID] = teacher
	return nil
}

func (m *MockTeacherUseCases) RetrieveTeacher(ctx context.Context, id string) (*models.Teacher, error) {
	if teacher, ok := m.Teachers[id]; ok {
		return teacher, nil
	}
	return nil, errors.New("teacher not found")
}

func (m *MockTeacherUseCases) RetrieveAllTeachers(ctx context.Context) ([]*models.Teacher, error) {
	var teachers []*models.Teacher
	for _, teacher := range m.Teachers {
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func (m *MockTeacherUseCases) ModifyTeacher(ctx context.Context, id string, teacher *models.Teacher) error {
	if _, ok := m.Teachers[id]; ok {
		teacher.HashedID = id
		m.Teachers[id] = teacher
		return nil
	}
	return errors.New("teacher not found")
}

func (m *MockTeacherUseCases) RemoveTeacher(ctx context.Context, id string) error {
	if _, ok := m.Teachers[id]; ok {
		delete(m.Teachers, id)
		return nil
	}
	return errors.New("teacher not found")
}

func TestAddTeacherHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockUseCases := NewMockTeacherUseCases()
	handlers := handlersfirebase.NewTeacherHandlers(mockUseCases)

	router.POST("/teacher", handlers.AddTeacher)

	newTeacher := models.Teacher{
		Name:      "pee dans",
		FirstName: "pee",
		LastName:  "dans",
		Subject:   "Math",
		Email:     "pee@gmail.com",
		Phone:     "0873234231",
	}
	body, _ := json.Marshal(newTeacher)
	req, _ := http.NewRequest(http.MethodPost, "/teacher", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRetrieveTeacherHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockUseCases := NewMockTeacherUseCases()
	handlers := handlersfirebase.NewTeacherHandlers(mockUseCases)

	// Add a teacher for retrieval
	teacher := &models.Teacher{
		Name:      "pee dans",
		FirstName: "pee",
		LastName:  "dans",
		Subject:   "Math",
		Email:     "pee@gmail.com",
		Phone:     "0873234231",
	}
	err := mockUseCases.AddTeacher(context.Background(), teacher)
	if err != nil {
		return
	}

	router.GET("/teacher/:id", handlers.RetrieveTeacher)
	req, _ := http.NewRequest(http.MethodGet, "/teacher/"+teacher.HashedID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRetrieveAllTeachersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockUseCases := NewMockTeacherUseCases()
	handlers := handlersfirebase.NewTeacherHandlers(mockUseCases)

	router.GET("/teachers", handlers.RetrieveAllTeachers)
	req, _ := http.NewRequest(http.MethodGet, "/teachers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestModifyTeacherHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockUseCases := NewMockTeacherUseCases()
	handlers := handlersfirebase.NewTeacherHandlers(mockUseCases)

	// Add a teacher to modify
	teacher := &models.Teacher{
		Name:      "pee dans",
		FirstName: "pee",
		LastName:  "dans",
		Subject:   "Math",
		Email:     "pee@gmail.com",
		Phone:     "0873234231",
	}
	err := mockUseCases.AddTeacher(context.Background(), teacher)
	if err != nil {
		return
	}

	updatedTeacher := models.Teacher{
		Name:      "jack queen",
		FirstName: "jack",
		LastName:  "queen",
		Subject:   "Tech",
		Email:     "jackking@gmail.com",
		Phone:     "1234567890",
	}
	body, _ := json.Marshal(updatedTeacher)
	router.PUT("/teacher/:id", handlers.ModifyTeacher)
	req, _ := http.NewRequest(http.MethodPut, "/teacher/"+teacher.HashedID, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTeacherHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockUseCases := NewMockTeacherUseCases()
	handlers := handlersfirebase.NewTeacherHandlers(mockUseCases)

	teacher := &models.Teacher{
		Name:      "jack queen",
		FirstName: "jack",
		LastName:  "queen",
		Subject:   "Tech",
		Email:     "jackking@gmail.com",
		Phone:     "1234567890",
	}
	err := mockUseCases.AddTeacher(context.Background(), teacher)
	if err != nil {
		return
	}

	router.DELETE("/teacher/:id", handlers.RemoveTeacher)
	req, _ := http.NewRequest(http.MethodDelete, "/teacher/"+teacher.HashedID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
