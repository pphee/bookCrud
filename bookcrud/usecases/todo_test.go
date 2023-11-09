package usecases

import (
	model "book/bookcrud"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindAll() ([]*model.Book, error) {
	args := m.Called()
	return args.Get(0).([]*model.Book), args.Error(1)
}

func (m *mockRepository) Update(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *mockRepository) New(book *model.Book) (*model.Book, error) {
	args := m.Called(book)
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *mockRepository) GetByID(id int) (*model.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *mockRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBookUsecase_CreateBook(t *testing.T) {
	mockRepo := new(mockRepository)

	book := &model.Book{
		Book:   "Sample Text",
		Author: "Author",
		Title:  "Title",
		ID:     1,
	}

	mockRepo.On("New", book).Return(book, nil)

	usecase := NewBookUsecase(mockRepo)
	err := usecase.CreateBook(book)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecase_GetAllBooks(t *testing.T) {
	mockRepo := new(mockRepository)
	var mockBooks []*model.Book
	mockBook := model.Book{ID: 1, Book: "Sample Text", Author: "Author", Title: "Title", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockBooks = append(mockBooks, &mockBook)

	mockRepo.On("FindAll").Return(mockBooks, nil)

	usecase := NewBookUsecase(mockRepo)
	books, err := usecase.GetAllBooks()

	assert.NoError(t, err)
	assert.NotNil(t, books)
	assert.Equal(t, 1, len(books))
	assert.Equal(t, "Sample Text", books[0].Book)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecase_UpdateBook(t *testing.T) {
	mockRepo := new(mockRepository)
	mockBook := &model.Book{ID: 1, Book: "Updated Text", Author: "Author", Title: "Title"}

	mockRepo.On("Update", mockBook).Return(nil)

	usecase := NewBookUsecase(mockRepo)
	err := usecase.UpdateBook(mockBook)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecase_GetBookByID(t *testing.T) {
	mockRepo := new(mockRepository)
	mockBook := &model.Book{ID: 1, Book: "Sample Text", Author: "Author", Title: "Title"}

	mockRepo.On("GetByID", 1).Return(mockBook, nil)

	usecase := NewBookUsecase(mockRepo)
	book, err := usecase.GetBookByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, mockBook, book)
	mockRepo.AssertExpectations(t)
}

func TestBookUsecase_DeleteBook(t *testing.T) {
	mockRepo := new(mockRepository)

	mockRepo.On("Delete", 1).Return(nil)

	usecase := NewBookUsecase(mockRepo)
	err := usecase.DeleteBook(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
