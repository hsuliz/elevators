package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hsuliz/elevators/internal/domain"
)

type SystemHandler struct {
	domainSystem *domain.System
	clients      map[*websocket.Conn]bool
	mu           sync.RWMutex
}

func NewSystemHandler(domainSystem *domain.System) *SystemHandler {
	return &SystemHandler{
		domainSystem: domainSystem,
		clients:      make(map[*websocket.Conn]bool),
	}
}

func (s *SystemHandler) CallElevator(c *gin.Context) {
	floorParam := c.Param("floor")
	floorNumber, err := strconv.Atoi(floorParam)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid floor number",
		})
		return
	}

	s.domainSystem.Call(floorNumber)
	c.Status(http.StatusOK)
}

func (s *SystemHandler) Activity(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	s.clients[conn] = true

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			delete(s.clients, conn)
			break
		}
	}
}

func (s *SystemHandler) ProcessActivity() {
	for activity := range s.domainSystem.ActivityCh {
		for client := range s.clients {
			if err := client.WriteJSON(&activity); err != nil {
				log.Printf("[Server] Error: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now (not recommended for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
