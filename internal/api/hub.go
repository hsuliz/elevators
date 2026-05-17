package api

import (
	"sync"

	"github.com/gorilla/websocket"
)

type client struct {
	conn   *websocket.Conn
	sendCh chan []byte
}

func newClient(conn *websocket.Conn) *client {
	c := &client{
		conn:   conn,
		sendCh: make(chan []byte, 32),
	}
	go c.writePump()
	return c
}

func (c *client) writePump() {
	for payload := range c.sendCh {
		if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return
		}
	}
}

func (c *client) close() {
	close(c.sendCh)
	c.conn.Close()
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]*client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]*client),
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = newClient(conn)
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if c, ok := h.clients[conn]; ok {
		c.close()
		delete(h.clients, conn)
	}
}

func (h *Hub) Broadcast(payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		select {
		case c.sendCh <- payload:
		default:
			// client too slow; drop frame
		}
	}
}
