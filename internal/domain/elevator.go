package domain

import (
	"sync"
	"time"

	"github.com/hsuliz/elevators/internal/domain/types"
)

type Elevator struct {
	ID                int
	CurrentFloor      int
	DestinationFloors []int
	Status            types.Status

	mu       sync.Mutex
	requests chan int
	updateCh chan int
}

func NewElevator(id int) *Elevator {
	return &Elevator{
		ID:                id,
		DestinationFloors: make([]int, 0),
		Status:            types.IDLE,
		requests:          make(chan int),
		updateCh:          make(chan int),
	}
}

func (e *Elevator) Run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case floor := <-e.requests:
			e.addFloor(floor)
		case <-ticker.C:
			e.step()
		}
	}
}

func (e *Elevator) addFloor(floor int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.DestinationFloors = append(e.DestinationFloors, floor)

	select {
	case e.updateCh <- e.ID:
	default:
	}
}

func (e *Elevator) step() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.DestinationFloors) == 0 {
		e.Status = types.IDLE
		return
	}

	target := e.DestinationFloors[0]

	switch {
	default:
		e.Status = types.IDLE
		e.DestinationFloors = e.DestinationFloors[1:]
	case e.CurrentFloor < target:
		e.CurrentFloor++
		e.Status = types.UP
	case e.CurrentFloor > target:
		e.CurrentFloor--
		e.Status = types.DOWN
	}

	select {
	case e.updateCh <- e.ID:
	default:
	}
}
