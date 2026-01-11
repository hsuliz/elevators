package types

import (
	"github.com/hsuliz/elevators/internal/domain/types"
)

type ElevatorResponse struct {
	ID           int          `json:"id"`
	CurrentFloor int          `json:"currentFloor"`
	Status       types.Status `json:"status"`
}
