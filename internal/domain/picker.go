package domain

import (
	"math/rand"
)

type Picker interface {
	Pick(elevators []*Elevator) int
}

type NaivePicker struct{}

func NewNaivePicker() *NaivePicker { return &NaivePicker{} }

func (n NaivePicker) Pick(elevators []*Elevator) int {
	for i, e := range elevators {
		if e.Status == IDLE {
			return i
		}
	}
	return rand.Intn(len(elevators))
}
