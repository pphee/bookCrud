package usecasesfirebase

import (
	models "book/teachercrud"
	"book/teachercrud/repofirebase"
	"context"
)

type ITeacherUseCases interface {
	AddTeacher(ctx context.Context, teacher *models.Teacher) error
	RetrieveTeacher(ctx context.Context, id string) (*models.Teacher, error)
	RetrieveAllTeachers(ctx context.Context) ([]*models.Teacher, error)
	ModifyTeacher(ctx context.Context, id string, teacher *models.Teacher) error
	RemoveTeacher(ctx context.Context, id string) error
}

type TeacherUseCases struct {
	repo repofirebase.ITeacherRepository
}

func NewTeacherUseCases(repo repofirebase.ITeacherRepository) ITeacherUseCases {
	return &TeacherUseCases{
		repo: repo,
	}
}

func (uc *TeacherUseCases) AddTeacher(ctx context.Context, teacher *models.Teacher) error {
	return uc.repo.CreateTeacher(ctx, teacher)
}

func (uc *TeacherUseCases) RetrieveTeacher(ctx context.Context, id string) (*models.Teacher, error) {
	return uc.repo.GetTeacherByID(ctx, id)
}

func (uc *TeacherUseCases) RetrieveAllTeachers(ctx context.Context) ([]*models.Teacher, error) {
	return uc.repo.GetAllTeachers(ctx)
}

func (uc *TeacherUseCases) ModifyTeacher(ctx context.Context, id string, teacher *models.Teacher) error {
	return uc.repo.UpdateTeacher(ctx, id, teacher)
}

func (uc *TeacherUseCases) RemoveTeacher(ctx context.Context, id string) error {
	return uc.repo.DeleteTeacher(ctx, id)
}
