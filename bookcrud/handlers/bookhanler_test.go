package handlers

import (
	model "book/bookcrud"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBookUsecase struct {
	mock.Mock
}

func (m *mockBookUsecase) GetAllBooks() ([]*model.Book, error) {
	args := m.Called()
	return args.Get(0).([]*model.Book), args.Error(1)
}

func (m *mockBookUsecase) CreateBook(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *mockBookUsecase) DeleteBook(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockBookUsecase) GetBookByID(id int) (*model.Book, error) {
	args := m.Called(id)
	if item := args.Get(0); item != nil {
		return item.(*model.Book), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockBookUsecase) UpdateBook(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func TestListBooks(t *testing.T) {

	mockUsecase := new(mockBookUsecase)
	handler := NewBookHandler(mockUsecase)

	expectedBooks := []*model.Book{{ID: 1, Title: "Test Book", Author: "Test Author", Book: "New Content"}}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/books", handler.ListBooks)

	tests := []struct {
		description  string
		returnBooks  []*model.Book
		returnErr    error
		expectedCode int
		expectedBody string
	}{
		{
			description:  "Get all books",
			returnBooks:  expectedBooks,
			expectedCode: http.StatusOK,
			expectedBody: "", // You would fill this in with the expected JSON string of expectedBooks
		},
		{
			description:  "Internal server error",
			returnBooks:  nil,
			returnErr:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal server error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			// Reset expectations for the mock
			mockUsecase.ExpectedCalls = nil
			mockUsecase.Calls = nil

			mockUsecase.On("GetAllBooks").Return(tc.returnBooks, tc.returnErr).Once()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/books", nil)

			r.ServeHTTP(w, req)

			if tc.expectedCode == http.StatusOK {
				var booksResponse []*model.Book
				err := json.Unmarshal(w.Body.Bytes(), &booksResponse)
				assert.NoError(t, err)
				assert.Equal(t, expectedBooks, booksResponse)
			} else {
				assert.Contains(t, w.Body.String(), tc.expectedBody)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetBookByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mockBookUsecase)
	handler := NewBookHandler(mockUsecase)

	expectedBook := &model.Book{ID: 1, Title: "Test Book", Author: "Test Author", Book: "New Content"}

	r := gin.Default()
	r.GET("/books/:id", handler.GetByID)

	tests := []struct {
		description   string
		bookID        string
		mockUsecaseFn func()
		expectedCode  int
	}{
		{
			description: "should fetch a book successfully",
			bookID:      "1",
			mockUsecaseFn: func() {
				mockUsecase.On("GetBookByID", 1).Return(expectedBook, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			description:   "should return bad request for invalid ID",
			bookID:        "abc",
			mockUsecaseFn: func() {},
			expectedCode:  http.StatusBadRequest,
		},
		{
			description: "should return not found if the book does not exist",
			bookID:      "2",
			mockUsecaseFn: func() {
				mockUsecase.On("GetBookByID", 2).Return(nil, errors.New("book not found"))
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.mockUsecaseFn()

			req, err := http.NewRequest(http.MethodGet, "/books/"+tc.bookID, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedCode == http.StatusOK {
				var book model.Book
				err = json.NewDecoder(rr.Body).Decode(&book)
				assert.NoError(t, err)
				assert.Equal(t, expectedBook, &book)
			}
		})
	}

	// Check if all expectations were met
	mockUsecase.AssertExpectations(t)
}

func TestDeleteBook(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockUsecase := new(mockBookUsecase)
	h := NewBookHandler(mockUsecase)
	r.DELETE("/book/:id", h.DeleteBook)

	// Define the test cases
	tests := []struct {
		description   string
		bookID        string
		mockUsecaseFn func()
		expectedCode  int
	}{
		{
			description: "should delete a book successfully",
			bookID:      "1",
			mockUsecaseFn: func() {
				mockUsecase.On("DeleteBook", 1).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			description:   "should return bad request for invalid ID",
			bookID:        "abc",
			mockUsecaseFn: func() {},
			expectedCode:  http.StatusBadRequest,
		},
		{
			description: "should return internal server error if usecase fails",
			bookID:      "2",
			mockUsecaseFn: func() {
				mockUsecase.On("DeleteBook", 2).Return(errors.New("some error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.mockUsecaseFn()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/book/%s", tt.bookID), nil)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedCode, resp.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestCreateBook(t *testing.T) {
	// Arrange
	mockUsecase := new(mockBookUsecase)
	handler := NewBookHandler(mockUsecase)

	inputBook := &model.Book{Title: "New Book", Author: "New Author", Book: "New Content"}
	invalidInputBook := &model.Book{Title: "", Author: "", Book: "New Content"} // Invalid because the fields are empty

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/books", handler.CreateBook)

	tests := []struct {
		description  string
		inputBook    *model.Book
		prepare      func()
		expectedCode int
	}{
		{
			description: "create book success",
			inputBook:   inputBook,
			prepare: func() {
				mockUsecase.On("CreateBook", mock.AnythingOfType("*model.Book")).Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			description: "create book bad request",
			inputBook:   invalidInputBook,
			prepare: func() {
				// No need to mock usecase for invalid input as it should fail before reaching usecase logic.
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			description: "create book internal server error",
			inputBook:   inputBook,
			prepare: func() {
				mockUsecase.On("CreateBook", mock.AnythingOfType("*model.Book")).Return(errors.New("internal server error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			mockUsecase.ExpectedCalls = nil
			mockUsecase.Calls = nil

			tt.prepare()

			body, _ := json.Marshal(tt.inputBook)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusCreated {
				var bookResponse model.Book
				err := json.Unmarshal(w.Body.Bytes(), &bookResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.inputBook.Title, bookResponse.Title)
			}
		})
	}
}

func TestBookUpdateBook(t *testing.T) {
	mockUsecase := new(mockBookUsecase)
	h := NewBookHandler(mockUsecase)

	gin.SetMode(gin.TestMode)

	tests := []struct {
		description   string
		bookID        string
		requestBody   *model.Book
		mockUsecaseFn func()
		expectedCode  int
		want          string
	}{
		{
			description: "success",
			bookID:      "1",
			requestBody: &model.Book{Book: "Updated Text", Author: "Author", Title: "Title"},
			mockUsecaseFn: func() {
				mockBook := &model.Book{ID: 1, Book: "Updated Text", Author: "Author", Title: "Title"}
				mockUsecase.On("GetBookByID", mockBook.ID).Return(mockBook, nil)
				mockUsecase.On("UpdateBook", mock.Anything).Return(nil)
			},
			expectedCode: http.StatusOK,
			want:         `{"message":"Book updated successfully"}`,
		},
		{
			description: "error",
			bookID:      "2",
			mockUsecaseFn: func() {
				mockBook := &model.Book{ID: 2, Book: "Updated Text", Author: "Author", Title: "Title"}
				mockUsecase.On("GetBookByID", mockBook.ID).Return(nil, errors.New("error"))
				mockUsecase.On("UpdateBook", mock.Anything).Return(errors.New("error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			description: "bad request",
			bookID:      "abc",
			mockUsecaseFn: func() {
				mockBook := &model.Book{ID: 1, Book: "Updated Text", Author: "Author", Title: "Title"}
				mockUsecase.On("GetBookByID", mockBook.ID).Return(mockBook, nil)
				mockUsecase.On("UpdateBook", mock.Anything).Return(nil)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			r := gin.Default()
			r.PUT("/books/:id", h.UpdateBook)

			tt.mockUsecaseFn()

			mockBook := &model.Book{ID: 1, Book: "Updated Text", Author: "Author", Title: "Title"}
			body, _ := json.Marshal(mockBook)
			req, _ := http.NewRequest(http.MethodPut, "/books/"+tt.bookID, bytes.NewBuffer(body))
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedCode, resp.Code)

			if resp.Code == http.StatusOK {
				assert.JSONEq(t, tt.want, resp.Body.String())
			}

			mockUsecase.AssertExpectations(t)
		})
	}

}
