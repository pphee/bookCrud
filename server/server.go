package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IServer interface {
	StartGin()
	GetServer() *server
}

type server struct {
	db  *gorm.DB
	app *gin.Engine
}

func NewServer(db *gorm.DB) IServer {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	return &server{
		app: app,
		db:  db,
	}
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) StartGin() {
	api := s.app.Group("/api")
	modules := InitModule(api, s)
	modules.BookModule()
	port := ":8080"
	log.Printf("Server is starting on %v", port)

	// Start the Gin
	if err := s.app.Run(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}

}
