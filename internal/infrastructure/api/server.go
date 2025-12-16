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

	router.POST("/call/:floor", systemHandler.CallElevator)
	router.GET("/ws", systemHandler.Activity)
	go systemHandler.ProcessActivity()

	return &APIServer{ginEngine: router}
}

func (s APIServer) Start(addr string) {
	log.Fatal(s.ginEngine.Run(addr))
}
