package handlers

import (
	model "book/bookcrud"
	"book/bookcrud/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IBookHandler interface {
	ListBooks(c *gin.Context)
	CreateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
	GetByID(c *gin.Context)
	UpdateBook(c *gin.Context)
}

type bookHandler struct {
	usecase usecases.IBookUsecase
}

func NewBookHandler(usecase usecases.IBookUsecase) IBookHandler {
	return &bookHandler{
		usecase: usecase,
	}
}

func (h *bookHandler) ListBooks(c *gin.Context) {
	books, err := h.usecase.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &books)
}

func (h *bookHandler) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &book)
}

func (h *bookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.usecase.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *bookHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	book, err := h.usecase.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &book)
}

func (h *bookHandler) UpdateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid beer ID"})
		return
	}

	bookResponse, err := h.usecase.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch beer details"})
		return
	}

	bookResponse.Book = book.Book
	bookResponse.Author = book.Author
	bookResponse.Title = book.Title

	if err := h.usecase.UpdateBook(bookResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}
