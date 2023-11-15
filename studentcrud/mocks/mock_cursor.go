package mocks

import (
	models "book/studentcrud" // Replace with your models' package path
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
)

type MockCursor struct {
	mock.Mock
	Data    []*models.Student
	Current int
}

func (m *MockCursor) Next(ctx context.Context) bool {
	m.Called(ctx)
	if m.Current < len(m.Data) {
		m.Current++
		return true
	}
	return false
}

func (m *MockCursor) Decode(val interface{}) error {
	m.Called(val)
	if m.Current <= len(m.Data) {
		*val.(*models.Student) = *m.Data[m.Current-1]
		return nil
	}
	return errors.New("out of range")
}

func (m *MockCursor) Close(ctx context.Context) error {
	m.Called(ctx)
	return nil
}
