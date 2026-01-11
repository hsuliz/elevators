package picker

import (
	"math/rand"

	"github.com/hsuliz/elevators/internal/domain/types"
)

type Naive struct{}

func NewNaive() *Naive { return &Naive{} }

func (n Naive) Pick(elevators []*types.Elevator) int {
	for i, e := range elevators {
		if e.Status == types.IDLE {
			return i
		}
	}
	return rand.Intn(len(elevators))
}
