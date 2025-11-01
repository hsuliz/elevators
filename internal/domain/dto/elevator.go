package dto

import "sync"

type Elevator struct {
	CurrentFlor int
	Status      Status
	Locker      sync.Locker
}

func NewElevator() *Elevator {
	return &Elevator{Locker: &sync.Mutex{}}
}
