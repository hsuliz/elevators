package domain

import (
	"math/rand"

	"github.com/hsuliz/elevators/internal/domain/dto"
)

type Naive struct{}

func NewNaive() *Naive { return &Naive{} }

func (n Naive) Pick(elevators []*dto.Elevator) int {
	for i, e := range elevators {
		if e.Status != dto.IDLE {
			return i
		}
	}
	return rand.Intn(len(elevators))
}
