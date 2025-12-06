package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hsuliz/elevators/internal/domain"
	"github.com/hsuliz/elevators/internal/port/api/types"
)

type SystemHandler struct {
	domainSystem *domain.System
}

func NewSystemHandler(domainSystem *domain.System) *SystemHandler {
	return &SystemHandler{domainSystem}
}

func (s SystemHandler) CallElevator(c *gin.Context) {
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

func (s SystemHandler) SystemStatus(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	for activity := range s.domainSystem.Activity() {
		statusRes := types.ElevatorResponse{
			ID: activity.ID, CurrentFloor: activity.CurrentFloor, Status: activity.Status,
		}
		data, err := json.Marshal(statusRes)
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now (not recommended for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
