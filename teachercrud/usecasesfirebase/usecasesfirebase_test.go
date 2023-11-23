package usecasesfirebase_test

import (
	models "book/teachercrud"
	"book/teachercrud/usecasesfirebase"
	"context"
	"strconv"
	"testing"
)

type MockTeacherRepository struct {
	Teachers map[string]*models.Teacher
}

func (m *MockTeacherRepository) CreateTeacher(ctx context.Context, teacher *models.Teacher) error {
	m.Teachers[teacher.HashedID] = teacher
	return nil
}

func (m *MockTeacherRepository) GetTeacherByID(ctx context.Context, id string) (*models.Teacher, error) {
	if teacher, exists := m.Teachers[id]; exists {
		return teacher, nil
	}
	return nil, nil
}

func (m *MockTeacherRepository) GetAllTeachers(ctx context.Context) ([]*models.Teacher, error) {
	var allTeachers []*models.Teacher
	for _, teacher := range m.Teachers {
		allTeachers = append(allTeachers, teacher)
	}
	return allTeachers, nil
}

func (m *MockTeacherRepository) UpdateTeacher(ctx context.Context, id string, teacher *models.Teacher) error {
	if _, exists := m.Teachers[id]; exists {
		m.Teachers[id] = teacher
		return nil
	}
	return nil
}

func (m *MockTeacherRepository) DeleteTeacher(ctx context.Context, id string) error {
	delete(m.Teachers, id)
	return nil
}

func TestAddTeacher(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockTeacherRepository{Teachers: make(map[string]*models.Teacher)}
	uc := usecasesfirebase.NewTeacherUseCases(mockRepo)

	teacher := &models.Teacher{
		Name:      "pee dans",
		FirstName: "pee",
		LastName:  "dans",
		Subject:   "Math",
		Email:     "pee@gmail.com",
		Phone:     "1234567890",
	}
	teacher.HashedID = teacher.HashID()

	err := uc.AddTeacher(ctx, teacher)
	if err != nil {
		t.Errorf("Failed to add teacher: %s", err)
	}

	if _, exists := mockRepo.Teachers[teacher.HashedID]; !exists {
		t.Errorf("Teacher was not added correctly")
	}
}

func TestGetTeacherByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockTeacherRepository{Teachers: make(map[string]*models.Teacher)}
	uc := usecasesfirebase.NewTeacherUseCases(mockRepo)

	expectedTeacher := &models.Teacher{
		Name:      "jack",
		FirstName: "queen",
		LastName:  "king",
		Subject:   "Math",
		Email:     "pee@gmail.com",
		Phone:     "1234567890",
	}
	expectedTeacher.HashedID = expectedTeacher.HashID()
	mockRepo.Teachers[expectedTeacher.HashedID] = expectedTeacher

	teacher, err := uc.RetrieveTeacher(ctx, expectedTeacher.HashedID)
	if err != nil {
		t.Errorf("Failed to retrieve teacher: %s", err)
	}
	if teacher == nil || teacher.HashedID != expectedTeacher.HashedID {
		t.Errorf("Retrieved teacher does not match expected")
	}
}

func TestGetAllTeachers(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockTeacherRepository{Teachers: make(map[string]*models.Teacher)}
	uc := usecasesfirebase.NewTeacherUseCases(mockRepo)

	for i := 0; i < 3; i++ {
		teacher := &models.Teacher{
			Name:      "mother" + strconv.Itoa(i),
			FirstName: "father" + strconv.Itoa(i),
			LastName:  "motherfather" + strconv.Itoa(i),
			Subject:   "computer" + strconv.Itoa(i),
			Email:     "teacher" + strconv.Itoa(i) + "@gmail.com",
			Phone:     "1234567890",
		}
		teacher.HashedID = teacher.HashID()
		mockRepo.Teachers[teacher.HashedID] = teacher
	}

	teachers, err := uc.RetrieveAllTeachers(ctx)
	if err != nil {
		t.Errorf("Failed to retrieve all teachers: %s", err)
	}
	if len(teachers) != 3 {
		t.Errorf("Expected 3 teachers, got %d", len(teachers))
	}
}

func TestUpdateTeacher(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockTeacherRepository{Teachers: make(map[string]*models.Teacher)}
	uc := usecasesfirebase.NewTeacherUseCases(mockRepo)

	teacher := &models.Teacher{
		Name:      "james",
		FirstName: "bank",
		LastName:  "jamesbank",
		Subject:   "sport",
		Email:     "jack@gmail.com",
		Phone:     "1234567890",
	}
	teacher.HashedID = teacher.HashID()
	mockRepo.Teachers[teacher.HashedID] = teacher

	updatedTeacher := &models.Teacher{
		Name:      "iphone",
		FirstName: "samsung",
		LastName:  "iphone samsung",
		Subject:   "Technology",
		Email:     "iphone@gmail.com",
		Phone:     "1234567890",
	}

	err := uc.ModifyTeacher(ctx, teacher.HashedID, updatedTeacher)
	if err != nil {
		t.Errorf("Failed to update teacher: %s", err)
	}
	if mockRepo.Teachers[teacher.HashedID].Name != "Updated Name" {
		t.Errorf("Teacher was not updated correctly")
	}
}

func TestDeleteTeacher(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockTeacherRepository{Teachers: make(map[string]*models.Teacher)}
	uc := usecasesfirebase.NewTeacherUseCases(mockRepo)

	teacher := &models.Teacher{
		Name:      "five",
		FirstName: "four",
		LastName:  "fivefour",
		Subject:   "five",
		Email:     "pee@gmail.com",
		Phone:     "1234567890",
	}
	teacher.HashedID = teacher.HashID()
	mockRepo.Teachers[teacher.HashedID] = teacher

	err := uc.RemoveTeacher(ctx, teacher.HashedID)
	if err != nil {
		t.Errorf("Failed to delete teacher: %s", err)
	}
	if _, exists := mockRepo.Teachers[teacher.HashedID]; exists {
		t.Errorf("Teacher was not deleted correctly")
	}
}
