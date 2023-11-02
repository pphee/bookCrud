package usecases

import (
	model "book/bookcrud"
	"book/bookcrud/repo"
)

type IBookUsecase interface {
	CreateBook(book *model.Book) error
	GetAllBooks() ([]*model.Book, error)
	UpdateBook(book *model.Book) error
	GetBookByID(id int) (*model.Book, error)
	DeleteBook(id int) error
}

type BookUsecase struct {
	repo repo.IBookRepository
}

func NewBookUsecase(repo repo.IBookRepository) IBookUsecase {
	return &BookUsecase{repo: repo}
}

func (u *BookUsecase) CreateBook(book *model.Book) error {
	return u.repo.New(book)
}

func (u *BookUsecase) GetAllBooks() ([]*model.Book, error) {
	return u.repo.FindAll()
}

func (u *BookUsecase) UpdateBook(book *model.Book) error {
	return u.repo.Update(book)
}

func (u *BookUsecase) GetBookByID(id int) (*model.Book, error) {
	return u.repo.GetByID(id)
}

func (u *BookUsecase) DeleteBook(id int) error {
	return u.repo.Delete(id)
}
