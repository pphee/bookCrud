package generator

import (
	models "book/studentcrud"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

func GenerateStudents() ([]models.Student, error) {
	var s []models.Student
	err := faker.FakeData(
		&s,
		options.WithRandomMapAndSliceMaxSize(100),
		options.WithRandomMapAndSliceMinSize(80),
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}
