package repofirebase

import (
	models "book/teachercrud"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
)

type ITeacherRepository interface {
	CreateTeacher(ctx context.Context, teacher *models.Teacher) error
	GetTeacherByID(ctx context.Context, id string) (*models.Teacher, error)
	GetAllTeachers(ctx context.Context) ([]*models.Teacher, error)
	UpdateTeacher(ctx context.Context, id string, teacher *models.Teacher) error
	DeleteTeacher(ctx context.Context, id string) error
}

type TeacherRepository struct {
	collection *firestore.CollectionRef
}

func NewTeacherRepository(collection *firestore.CollectionRef) ITeacherRepository {
	return &TeacherRepository{collection: collection}
}

func (r *TeacherRepository) CreateTeacher(ctx context.Context, teacher *models.Teacher) error {
	_, err := r.collection.Doc(teacher.ID).Set(ctx, teacher)
	return err
}

func (r *TeacherRepository) GetTeacherByID(ctx context.Context, id string) (*models.Teacher, error) {
	doc, err := r.collection.Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var teacher models.Teacher
	if err := doc.DataTo(&teacher); err != nil {
		return nil, err
	}
	teacher.ID = doc.Ref.ID
	return &teacher, nil
}

func (r *TeacherRepository) GetAllTeachers(ctx context.Context) ([]*models.Teacher, error) {
	var teachers []*models.Teacher

	iter := r.collection.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var teacher models.Teacher
		if err := doc.DataTo(&teacher); err != nil {
			return nil, err
		}
		teacher.ID = doc.Ref.ID
		teachers = append(teachers, &teacher)
	}

	return teachers, nil
}

func (r *TeacherRepository) UpdateTeacher(ctx context.Context, id string, teacher *models.Teacher) error {
	_, err := r.collection.Doc(id).Set(ctx, teacher, firestore.MergeAll)
	return err
}

func (r *TeacherRepository) DeleteTeacher(ctx context.Context, id string) error {
	_, err := r.collection.Doc(id).Delete(ctx)
	return err
}
