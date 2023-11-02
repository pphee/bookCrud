package server

import (
	"book/bookcrud/handlers"
	"book/bookcrud/repo"
	"book/bookcrud/usecases"

	"github.com/gin-gonic/gin"
)

type IModuleFactory interface {
	bookModule()
}

type moduleFactory struct {
	r *gin.RouterGroup
	s *server
}

func InitModule(r *gin.RouterGroup, s *server) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
	}
}

func (mf *moduleFactory) bookModule() {
	gormStore := repo.NewGormStore(mf.s.db)
	useCases := usecases.NewBookUsecase(gormStore)
	bookHandlers := handlers.NewBookHandler(useCases)
	bookRouter := mf.r.Group("/books")
	bookRouter.GET("/:id", bookHandlers.GetByID)
	bookRouter.GET("", bookHandlers.ListBooks)
	bookRouter.POST("", bookHandlers.CreateBook)
	bookRouter.PUT("/:id", bookHandlers.UpdateBook)
	bookRouter.DELETE("/:id", bookHandlers.DeleteBook)
}
