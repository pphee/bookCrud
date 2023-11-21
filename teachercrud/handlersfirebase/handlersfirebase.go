package handlersfirebase

import (
	models "book/teachercrud"
	usecases "book/teachercrud/usecasesfirebase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeacherHandlers struct {
	Usecases usecases.ITeacherUseCases
}

func NewTeacherHandlers(usecases usecases.ITeacherUseCases) *TeacherHandlers {
	return &TeacherHandlers{
		Usecases: usecases,
	}
}

func (h *TeacherHandlers) AddTeacher(c *gin.Context) {
	var teacher models.Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Usecases.AddTeacher(c.Request.Context(), &teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teacher)
}

func (h *TeacherHandlers) RetrieveTeacher(c *gin.Context) {
	id := c.Param("id")
	teacher, err := h.Usecases.RetrieveTeacher(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

func (h *TeacherHandlers) RetrieveAllTeachers(c *gin.Context) {
	teachers, err := h.Usecases.RetrieveAllTeachers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

func (h *TeacherHandlers) ModifyTeacher(c *gin.Context) {
	id := c.Param("id")
	var teacher models.Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Usecases.ModifyTeacher(c.Request.Context(), id, &teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher updated successfully"})
}

func (h *TeacherHandlers) RemoveTeacher(c *gin.Context) {
	id := c.Param("id")
	if err := h.Usecases.RemoveTeacher(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher removed successfully"})
}
