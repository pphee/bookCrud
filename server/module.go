package server

import (
	"book/bookcrud/handlers"
	"book/bookcrud/repo"
	"book/bookcrud/usecases"

	"github.com/gin-gonic/gin"
)

type IModuleFactory interface {
	BookModule()
	Repository() repo.IBookRepository
	Usecase() usecases.IBookUsecase
	Handler() handlers.IBookHandler
}

type moduleFactory struct {
	r          *gin.RouterGroup
	s          *server
	repository repo.IBookRepository
	usecase    usecases.IBookUsecase
	handler    handlers.IBookHandler
}

func InitModule(r *gin.RouterGroup, s *server) IModuleFactory {
	mf := &moduleFactory{
		r: r,
		s: s,
	}
	bookRepository := repo.NewGormStore(s.db)
	bookUsecase := usecases.NewBookUsecase(bookRepository)
	bookHandler := handlers.NewBookHandler(bookUsecase)

	mf.repository = bookRepository
	mf.usecase = bookUsecase
	mf.handler = bookHandler

	return mf
}

func (mf *moduleFactory) BookModule() {

	bookRouter := mf.r.Group("/books")

	bookRouter.GET("/:id", mf.handler.GetByID)
	bookRouter.GET("", mf.handler.ListBooks)
	bookRouter.POST("", mf.handler.CreateBook)
	bookRouter.PUT("/:id", mf.handler.UpdateBook)
	bookRouter.DELETE("/:id", mf.handler.DeleteBook)
}

func (f *moduleFactory) Repository() repo.IBookRepository { return f.repository }
func (f *moduleFactory) Usecase() usecases.IBookUsecase   { return f.usecase }
func (f *moduleFactory) Handler() handlers.IBookHandler   { return f.handler }
