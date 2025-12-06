package main

import (
	"github.com/hsuliz/elevators/internal/domain"
	"github.com/hsuliz/elevators/internal/port/api"
	"github.com/hsuliz/elevators/internal/port/api/handler"
)

func main() {
	elevator1 := domain.NewElevator(1)
	elevator2 := domain.NewElevator(2)
	elevators := []*domain.Elevator{elevator1, elevator2}
	naivePicker := domain.NewNaive()
	system := domain.NewSystem(elevators, naivePicker, 11)

	systemHandler := handler.NewSystem(system)
	webSocketHandler := handler.NewWebSocket(system)
	api.NewServer(systemHandler, webSocketHandler).Start(":8080")
}
