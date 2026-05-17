package api

import (
	"context"
	"encoding/json"
	"github.com/hsuliz/elevators/internal/domain"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type ElevatorDTO struct {
	ID           int           `json:"id"`
	CurrentFloor int           `json:"currentFloor"`
	Status       domain.Status `json:"status"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins in dev
}

type Server struct {
	system *domain.System
	hub    *Hub
	mux    *http.ServeMux
}

func New(system *domain.System) *Server {
	s := &Server{
		system: system,
		hub:    NewHub(),
		mux:    http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/elevators", s.handleElevators)
	s.mux.HandleFunc("/floors", s.handleFloors)
	s.mux.HandleFunc("/call/", s.handleCall) // /call/{floor}
	s.mux.HandleFunc("/ws", s.handleWS)

	fs := http.FileServer(http.Dir("static"))
	s.mux.Handle("/", fs)
}

func (s *Server) handleElevators(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dtos := make([]ElevatorDTO, len(s.system.Elevators))
	for i, e := range s.system.Elevators {
		a := e.GetActivity()
		dtos[i] = ElevatorDTO{ID: a.ID, CurrentFloor: a.CurrentFloor, Status: a.Status}
	}

	jsonResponse(w, dtos)
}

func (s *Server) handleFloors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	floors := make(map[int]bool, s.system.FloorCount+1)
	for i := 0; i <= s.system.FloorCount; i++ {
		floors[i] = true
	}

	jsonResponse(w, floors)
}

func (s *Server) handleCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract floor from path: /call/3 → "3"
	part := strings.TrimPrefix(r.URL.Path, "/call/")
	floor, err := strconv.Atoi(part)
	if err != nil {
		http.Error(w, "invalid floor number", http.StatusBadRequest)
		return
	}

	if err := s.system.Call(floor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}
	s.hub.Register(conn)
	log.Printf("ws: client connected (%s)", conn.RemoteAddr())

	// read loop running so we detect client disconnect.
	go func() {
		defer func() {
			s.hub.Unregister(conn)
			conn.Close()
			log.Printf("ws: client disconnected (%s)", conn.RemoteAddr())
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}

func (s *Server) WatchAndBroadcast(ctx context.Context) {
	for _, e := range s.system.Elevators {
		go func(elev *domain.Elevator) {
			for {
				select {
				case <-elev.GetUpdates():
					a := elev.GetActivity()
					dto := ElevatorDTO{ID: a.ID, CurrentFloor: a.CurrentFloor, Status: a.Status}
					payload, err := json.Marshal(dto)
					if err != nil {
						log.Printf("marshal error: %v", err)
						continue
					}
					s.hub.Broadcast(payload)
				case <-ctx.Done():
					return
				}
			}
		}(e)
	}
}

func jsonResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
