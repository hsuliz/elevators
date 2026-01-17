package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hsuliz/elevators/internal/infrastructure/api/handler"
)

type APIServer struct {
	ginEngine *gin.Engine
}

func NewServer(systemHandler *handler.System) *APIServer {
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.File("./static" + c.Request.URL.Path)
	})

	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.POST("/call/:floor", systemHandler.CallElevator)

	router.GET("/ws", systemHandler.Activity)
	go systemHandler.ProcessActivity()

	router.GET("/elevators", systemHandler.GetElevators)

	return &APIServer{ginEngine: router}
}

func (s APIServer) Start(addr string) {
	log.Fatal(s.ginEngine.Run(addr))
}
