package handlersmongo

import (
	models "book/studentcrud"
	"book/studentcrud/usecasesmongo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IStudentHandler interface {
	PostStudent(c *gin.Context)
	GetStudent(c *gin.Context)
	GetAllStudents(c *gin.Context)
	UpdateStudent(c *gin.Context)
	DeleteStudent(c *gin.Context)
}

type StudentHandler struct {
	useCase usecasesmongo.IStudentUseCase
}

func NewStudentHandler(useCase usecasesmongo.IStudentUseCase) IStudentHandler {
	return &StudentHandler{
		useCase: useCase,
	}
}

func (h *StudentHandler) PostStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if student.FirstName == "" || student.LastName == "" || student.Age <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, missing fields"})
		return
	}

	id, err := h.useCase.CreateStudent(c.Request.Context(), &student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	id := c.Param("id")
	student, err := h.useCase.GetStudentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	students, err := h.useCase.GetAllStudents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if student.FirstName == "" || student.LastName == "" || student.Age <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, missing fields"})
		return
	}

	err := h.useCase.UpdateStudent(c.Request.Context(), id, student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.DeleteStudent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found or unable to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
