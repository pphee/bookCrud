package server

import (
	"book/bookcrud/handlers"
	"book/bookcrud/repo"
	"book/bookcrud/usecases"
	"book/studentcrud/handlersmongo"
	"book/studentcrud/repomongo"
	"book/studentcrud/usecasesmongo"
	"github.com/gin-gonic/gin"
)

type IModuleFactory interface {
	BookModule()
	StudentModule()
}

type moduleFactory struct {
	r *gin.RouterGroup
	s *server
}

func InitModule(r *gin.RouterGroup, s *server) IModuleFactory {
	mf := &moduleFactory{
		r: r,
		s: s,
	}

	return mf
}

func (mf *moduleFactory) BookModule() {
	bookRepository := repo.NewGormStore(mf.s.db)
	bookUsecase := usecases.NewBookUsecase(bookRepository)
	bookHandler := handlers.NewBookHandler(bookUsecase)

	bookRouter := mf.r.Group("/books")

	bookRouter.GET("/:id", bookHandler.GetByID)
	bookRouter.GET("", bookHandler.ListBooks)
	bookRouter.POST("", bookHandler.CreateBook)
	bookRouter.PUT("/:id", bookHandler.UpdateBook)
	bookRouter.DELETE("/:id", bookHandler.DeleteBook)
}

func (mf *moduleFactory) StudentModule() {
	encryptionKey := []byte("your-encryption-key-here")

	studentRepository := repomongo.NewStudentRepository(mf.s.mongoCollection, encryptionKey)
	studentUsecase := usecasesmongo.NewStudentUseCase(studentRepository)
	studentHandler := handlersmongo.NewStudentHandler(studentUsecase)

	studentRouter := mf.r.Group("/students")

	studentRouter.GET("/:id", studentHandler.GetStudent)
	studentRouter.GET("", studentHandler.GetAllStudents)
	studentRouter.POST("", studentHandler.PostStudent)
	studentRouter.PUT("/:id", studentHandler.UpdateStudent)
	studentRouter.DELETE("/:id", studentHandler.DeleteStudent)
}
