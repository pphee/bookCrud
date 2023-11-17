package generator

import (
	model "book/bookcrud"
	models "book/studentcrud"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

func GenerateStudents() ([]models.Student, error) {
	var s []models.Student
	err := faker.FakeData(
		&s,
		options.WithRandomMapAndSliceMaxSize(1),
		options.WithRandomMapAndSliceMinSize(1),
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GenerateBooks() ([]model.Book, error) {
	var b []model.Book
	err := faker.FakeData(
		&b,
		options.WithRandomMapAndSliceMaxSize(1),
		options.WithRandomMapAndSliceMinSize(1),
	)
	if err != nil {
		return nil, err
	}
	return b, nil
}
