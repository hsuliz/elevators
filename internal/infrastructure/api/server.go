package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hsuliz/elevators/internal/infrastructure/api/handler"
)

type Server struct {
	ginEngine *gin.Engine
}

func NewServer(systemHandler *handler.System) *Server {
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.File("./static" + c.Request.URL.Path)
	})

	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.POST("/call/:floor", systemHandler.CallElevator)
	router.GET("/ws", systemHandler.Activity)
	router.GET("/elevators", systemHandler.GetElevators)

	// Start background processors once (broadcast loop, etc.)
	systemHandler.Start()

	return &Server{ginEngine: router}
}

func (s *Server) Start(addr string) {
	log.Fatal(s.ginEngine.Run(addr))
}
