package domain

import (
	"math/rand"
)

type Naive struct{}

func NewNaive() *Naive { return &Naive{} }

func (n Naive) Pick(elevators []*Elevator) int {
	for i, e := range elevators {
		if e.Status == IDLE {
			return i
		}
	}
	return rand.Intn(len(elevators))
}
