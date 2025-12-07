package domain

import "sync"

type Elevator struct {
	ID                int
	CurrentFloor      int
	DestinationFloors []int
	Status            Status
	mu                sync.Mutex
}

func NewElevator(id int) *Elevator {
	return &Elevator{ID: id, DestinationFloors: make([]int, 0)}
}
