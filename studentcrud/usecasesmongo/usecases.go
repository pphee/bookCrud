package usecasesmongo

import (
	models "book/studentcrud"
	"book/studentcrud/repomongo"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StudentUseCase interface defines the business logic for Student entities.
type IStudentUseCase interface {
	CreateStudent(ctx context.Context, student *models.Student) (primitive.ObjectID, error)
	GetStudentByID(ctx context.Context, id string) (*models.Student, error)
	GetAllStudents(ctx context.Context) ([]*models.Student, error)
	UpdateStudent(ctx context.Context, id string, student models.Student) error
	DeleteStudent(ctx context.Context, id string) error
}

type studentUseCase struct {
	repo repomongo.IStudentRepository
}

func NewStudentUseCase(repo repomongo.IStudentRepository) IStudentUseCase {
	return &studentUseCase{
		repo: repo,
	}
}

func (uc *studentUseCase) CreateStudent(ctx context.Context, student *models.Student) (primitive.ObjectID, error) {
	insertResult, err := uc.repo.Create(ctx, student)
	if err != nil {
		return primitive.NilObjectID, err
	}
	id := insertResult.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (uc *studentUseCase) GetStudentByID(ctx context.Context, id string) (*models.Student, error) {
	student, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (uc *studentUseCase) GetAllStudents(ctx context.Context) ([]*models.Student, error) {
	students, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (uc *studentUseCase) UpdateStudent(ctx context.Context, id string, student models.Student) error {
	_, err := uc.repo.Update(ctx, id, &student)
	return err
}

func (uc *studentUseCase) DeleteStudent(ctx context.Context, id string) error {
	_, err := uc.repo.Delete(ctx, id)
	return err
}
