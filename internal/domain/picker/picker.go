package picker

import "github.com/hsuliz/elevators/internal/domain/types"

type Interface interface {
	Pick(elevators []*types.Elevator) int
}
