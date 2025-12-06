package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hsuliz/elevators/internal/port/api/handler"
)

type Server struct {
	ginEngine *gin.Engine
}

func NewServer(systemHandler *handler.System, webSocketHandler *handler.WebSocket) *Server {
	router := gin.Default()

	router.POST("/call/:floor", systemHandler.CallElevator)
	router.GET("/ws", webSocketHandler.HandleWebSocket)

	return &Server{ginEngine: router}
}

func (s Server) Start(addr string) {
	log.Fatal(s.ginEngine.Run(addr))
}
