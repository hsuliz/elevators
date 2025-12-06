package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	ginEngine *gin.Engine
}

func NewServer(systemHandler *SystemHandler) *APIServer {
	router := gin.Default()

	router.POST("/call/:floor", systemHandler.CallElevator)
	router.GET("/ws", systemHandler.SystemStatus)

	return &APIServer{ginEngine: router}
}

func (s APIServer) Start(addr string) {
	log.Fatal(s.ginEngine.Run(addr))
}
