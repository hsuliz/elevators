package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hsuliz/elevators/internal/domain"
	"github.com/hsuliz/elevators/internal/infrastructure/api/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now (not recommended for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type System struct {
	domainSystem *domain.System
	clients      map[*websocket.Conn]bool
	mu           sync.RWMutex
}

func NewSystem(domainSystem *domain.System) *System {
	systemHandler := &System{
		domainSystem: domainSystem,
		clients:      make(map[*websocket.Conn]bool),
	}
	return systemHandler
}

func (h *System) CallElevator(c *gin.Context) {
	floorParam := c.Param("floor")
	floorNumber, err := strconv.Atoi(floorParam)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid floor number",
		})
		return
	}

	h.domainSystem.Call(floorNumber)
	c.Status(http.StatusOK)
}

func (h *System) Activity(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			break
		}
	}
}

func (h *System) ProcessActivity() {
	for activity := range h.domainSystem.ActivityCh {
		h.mu.RLock()
		for client := range h.clients {
			activityRes := types.ElevatorResponse{
				ID:           activity.ID,
				CurrentFloor: activity.CurrentFloor,
				Status:       activity.Status,
			}
			if err := client.WriteJSON(activityRes); err != nil {
				client.Close()
				h.mu.RUnlock() // release before acquiring write lock
				h.mu.Lock()
				delete(h.clients, client)
				h.mu.Unlock()
				h.mu.RLock() // reacquire read lock
			}
		}
		h.mu.RUnlock()
	}
}

func (h *System) GetElevators(c *gin.Context) {
	responses := make([]types.ElevatorResponse, 0, len(h.domainSystem.Elevators))

	for _, elevator := range h.domainSystem.Elevators {
		responses = append(responses, types.ElevatorResponse{
			ID:           elevator.ID,
			CurrentFloor: elevator.CurrentFloor,
			Status:       elevator.Status,
		})
	}

	c.JSON(http.StatusOK, responses)
}
