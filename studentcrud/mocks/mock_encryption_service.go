package mocks

import "github.com/stretchr/testify/mock"

type MockEncryptionService struct {
	mock.Mock
}

func (m *MockEncryptionService) Encrypt(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockEncryptionService) Decrypt(data string) ([]byte, error) {
	args := m.Called(data)
	result, ok := args.Get(0).([]byte)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}
