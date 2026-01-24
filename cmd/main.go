package main

import (
	"github.com/hsuliz/elevators/internal/domain"
	"github.com/hsuliz/elevators/internal/infrastructure/api/handler"

	api "github.com/hsuliz/elevators/internal/infrastructure/api"
)

func main() {
	elevator1 := domain.NewElevator(1)
	elevator2 := domain.NewElevator(2)
	elevators := []*domain.Elevator{elevator1, elevator2}

	naivePicker := domain.NewNaivePicker()
	system := domain.NewSystem(elevators, naivePicker)

	systemHandler := handler.NewSystem(system)
	api.NewServer(systemHandler).Start(":8080")
}
