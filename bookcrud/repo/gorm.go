package repo

import (
	model "book/bookcrud"

	"gorm.io/gorm"
)

type IBookRepository interface {
	New(book *model.Book) error
	FindAll() ([]*model.Book, error)
	Update(book *model.Book) error
	GetByID(id int) (*model.Book, error)
	Delete(id int) error
}

type gormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) IBookRepository {
	return &gormStore{db: db}
}

func (s *gormStore) New(book *model.Book) error {
	return s.db.Create(book).Error
}

func (s *gormStore) FindAll() ([]*model.Book, error) {
	var books []*model.Book
	if result := s.db.Find(&books); result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func (s *gormStore) Update(book *model.Book) error {
	return s.db.Save(book).Error
}

func (s *gormStore) GetByID(id int) (*model.Book, error) {
	var book model.Book
	if result := s.db.First(&book, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func (s *gormStore) Delete(id int) error {
	return s.db.Delete(&model.Book{}, "id = ?", id).Error
}
