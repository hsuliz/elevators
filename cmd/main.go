package main

import (
	"github.com/hsuliz/elevators/internal/domain"

	api "github.com/hsuliz/elevators/internal/port/api"
)

func main() {
	elevator1 := domain.NewElevator(1)
	elevator2 := domain.NewElevator(2)
	elevators := []*domain.Elevator{elevator1, elevator2}
	naivePicker := domain.NewNaive()
	system := domain.NewSystem(elevators, naivePicker, 11)

	systemHandler := api.NewSystemHandler(system)
	api.NewServer(systemHandler).Start(":8080")
}
