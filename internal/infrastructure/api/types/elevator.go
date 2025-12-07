package types

import "github.com/hsuliz/elevators/internal/domain"

type ElevatorResponse struct {
	ID           int           `json:"id"`
	CurrentFloor int           `json:"currentFloor"`
	Status       domain.Status `json:"status"`
}
