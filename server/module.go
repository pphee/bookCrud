package server

import (
	"book/bookcrud/handlers"
	"book/bookcrud/repo"
	"book/bookcrud/usecases"
	"book/studentcrud/handlersmongo"
	"book/studentcrud/repomongo"
	"book/studentcrud/usecasesmongo"
	"book/teachercrud/handlersfirebase"
	"book/teachercrud/repofirebase"
	"book/teachercrud/usecasesfirebase"
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type IModuleFactory interface {
	BookModule()
	StudentModule()
	TeacherModule()
}

type moduleFactory struct {
	r              *gin.RouterGroup
	s              *server
	firebaseClient *firestore.Client
}

func InitModule(r *gin.RouterGroup, s *server, firebaseClient *firestore.Client) IModuleFactory {
	mf := &moduleFactory{
		r:              r,
		s:              s,
		firebaseClient: firebaseClient,
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
	encryptionKey := []byte("jack-queen-kingpok92deng")

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

func (mf *moduleFactory) TeacherModule() {

	teacherCollection := mf.firebaseClient.Collection("teachers")
	teacherRepo := repofirebase.NewTeacherRepository(teacherCollection)
	teacherUsecase := usecasesfirebase.NewTeacherUseCases(teacherRepo)
	teacherHandler := handlersfirebase.NewTeacherHandlers(teacherUsecase)

	teacherRouter := mf.r.Group("/teachers")

	teacherRouter.GET("/:id", teacherHandler.RetrieveTeacher)
	teacherRouter.GET("", teacherHandler.RetrieveAllTeachers)
	teacherRouter.POST("", teacherHandler.AddTeacher)
	teacherRouter.PUT("/:id", teacherHandler.ModifyTeacher)
	teacherRouter.DELETE("/:id", teacherHandler.RemoveTeacher)
}
