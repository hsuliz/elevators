package main

import (
	"github.com/hsuliz/elevators/internal/domain"
	"github.com/hsuliz/elevators/internal/domain/picker"
	"github.com/hsuliz/elevators/internal/domain/types"
	"github.com/hsuliz/elevators/internal/infrastructure/api/handler"

	api "github.com/hsuliz/elevators/internal/infrastructure/api"
)

func main() {
	elevator1 := types.NewElevator(1)
	elevator2 := types.NewElevator(2)
	elevators := []*types.Elevator{elevator1, elevator2}
	naivePicker := picker.NewNaive()
	system := domain.NewSystem(elevators, naivePicker, 11)

	systemHandler := handler.NewSystem(system)
	api.NewServer(systemHandler).Start(":8080")
}
