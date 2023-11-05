package repo

import (
	"log"
	"regexp"
	"testing"
	"time"

	model "book/bookcrud"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMock() (*gorm.DB, sqlmock.Sqlmock, IBookRepository) {
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatalf("failed to initialize sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initialize gorm: %v", err)
	}

	store := NewGormStore(gormDB)
	return gormDB, mock, store
}

func TestGormStore_New(t *testing.T) {
	_, mock, store := NewMock()
	mock.ExpectBegin()

	now := time.Now()

	want := &model.Book{
		Book:      "Sample Text",
		Author:    "Author",
		Title:     "Title",
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `books` (`book`,`author`,`title`,`created_at`,`updated_at`,`id`) VALUES (?,?,?,?,?,?)")).
		WithArgs(want.Book, want.Author, want.Title, sqlmock.AnyArg(), sqlmock.AnyArg(), want.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	got, err := store.New(&model.Book{
		Book:   "Sample Text",
		Author: "Author",
		Title:  "Title",
		ID:     1,
	})

	got.CreatedAt = want.CreatedAt
	got.UpdatedAt = want.UpdatedAt

	assert.NoError(t, err)
	assert.NotNil(t, want, got, "The created book does not match the expected book.")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormStore_FindAll(t *testing.T) {
	_, mock, store := NewMock()

	rows := sqlmock.NewRows([]string{"id", "book", "author", "title", "created_at", "updated_at"}).
		AddRow(1, "Sample Text", "Author", "Title", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books`")).WillReturnRows(rows)

	books, err := store.FindAll()

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormStore_Update(t *testing.T) {
	_, mock, store := NewMock()

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `books` SET `book`=?,`author`=?,`title`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")).
		WithArgs("Sample Text", "Author", "Title", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := store.Update(&model.Book{
		ID:     1,
		Book:   "Sample Text",
		Author: "Author",
		Title:  "Title",
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormStore_Delete(t *testing.T) {
	_, mock, store := NewMock()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `books` WHERE id = ?")).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := store.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGormStore_GetByID(t *testing.T) {
	_, mock, store := NewMock()

	want := &model.Book{
		ID:     1,
		Book:   "Sample Text",
		Author: "Author",
		Title:  "Title",
	}

	rows := sqlmock.NewRows([]string{"id", "book", "author", "title"}).
		AddRow(1, "Sample Text", "Author", "Title")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` WHERE id = ?")).
		WithArgs(1).
		WillReturnRows(rows)

	got, err := store.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, want, got, "Book retrieved by author does not match the expected result")
	assert.NoError(t, mock.ExpectationsWereMet())
}
