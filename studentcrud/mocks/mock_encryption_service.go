package mocks

import "github.com/stretchr/testify/mock"

type MockEncryptionService struct {
	mock.Mock
}

func (m *MockEncryptionService) Encrypt(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockEncryptionService) Decrypt(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}
