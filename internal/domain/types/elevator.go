package types

import "sync"

type Elevator struct {
	ID                int
	CurrentFloor      int
	DestinationFloors []int
	Status            Status
	Mu                sync.Mutex
}

func NewElevator(id int) *Elevator {
	return &Elevator{ID: id, DestinationFloors: make([]int, 0)}
}
