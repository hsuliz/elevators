package domain

import (
	"math"
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

type SmartPicker struct{}

func NewSmartPicker() *SmartPicker { return &SmartPicker{} }

func (s SmartPicker) Pick(elevators []*Elevator, targetFloor int) int {
	bestIdx := 0
	bestDist := math.MaxInt32
	bestIdle := false

	for i, e := range elevators {
		e.Lock()
		dist := abs(e.CurrentFloor - targetFloor)
		idle := e.Status == IDLE
		e.Unlock()

		better := false
		switch {
		case idle && !bestIdle:
			better = true
		case idle == bestIdle && dist < bestDist:
			better = true
		}

		if better {
			bestIdx = i
			bestDist = dist
			bestIdle = idle
		}
	}
	return bestIdx
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
