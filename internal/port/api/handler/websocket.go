package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/hsuliz/elevators/internal/domain"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now (not recommended for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocket struct {
	domainSystem *domain.System
}

func NewWebSocket(domainSystem *domain.System) *WebSocket {
	return &WebSocket{domainSystem: domainSystem}
}

func (w WebSocket) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	for status := range w.domainSystem.Status() {
		data, err := json.Marshal(status)
		if err != nil {
			log.Println(err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error writing message: %v", err)
			continue
		}
	}
}
