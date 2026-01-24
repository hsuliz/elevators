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
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type client struct {
	conn *websocket.Conn
	send chan types.ElevatorResponse
}

type hub struct {
	register   chan *client
	unregister chan *client
	broadcast  chan types.ElevatorResponse

	clients map[*client]struct{}
}

func newHub() *hub {
	return &hub{
		register:   make(chan *client, 16),
		unregister: make(chan *client, 16),
		broadcast:  make(chan types.ElevatorResponse, 128),
		clients:    make(map[*client]struct{}),
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				_ = c.conn.Close()
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
				}
			}
		}
	}
}

type System struct {
	domainSystem *domain.System

	hub  *hub
	once sync.Once
}

func NewSystem(domainSystem *domain.System) *System {
	return &System{
		domainSystem: domainSystem,
		hub:          newHub(),
	}
}

func (h *System) Start() {
	h.once.Do(func() {
		go h.hub.run()
		go h.processActivity()
	})
}

func (h *System) CallElevator(c *gin.Context) {
	floorParam := c.Param("floor")
	floorNumber, err := strconv.Atoi(floorParam)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid floor number"})
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

	cl := &client{
		conn: conn,
		send: make(chan types.ElevatorResponse, 32),
	}

	h.hub.register <- cl

	go func() {
		for msg := range cl.send {
			if err := cl.conn.WriteJSON(msg); err != nil {
				h.hub.unregister <- cl
				return
			}
		}
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.hub.unregister <- cl
			return
		}
	}
}

func (h *System) processActivity() {
	for activity := range h.domainSystem.ActivityCh {
		h.hub.broadcast <- types.ElevatorResponse{
			ID:           activity.ID,
			CurrentFloor: activity.CurrentFloor,
			Status:       activity.Status,
		}
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
