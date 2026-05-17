package domain

import (
	"log"
	"sync"
	"time"
)

type Status int

const (
	UP   = 1
	IDLE = 0
	DOWN = -1
)

type Activity struct {
	ID                int
	CurrentFloor      int
	DestinationFloors []int
	Status            Status
}

type Elevator struct {
	ID                int
	CurrentFloor      int
	DestinationFloors []int
	Status            Status

	requestCh chan int
	updateCh  chan struct{}
	mu        sync.Mutex
}

func NewElevator(id int) *Elevator {
	return &Elevator{
		ID:                id,
		DestinationFloors: make([]int, 0),
		Status:            IDLE,
	}
}

func (e *Elevator) TurnOn() {
	e.requestCh = make(chan int)
	e.updateCh = make(chan struct{}, 1)

	go e.requestor()
	go e.updator()

	log.Printf("elevator id %d: turned ON\n", e.ID)
}

func (e *Elevator) RequestFloor(floorNumber int) {
	log.Printf("elevator id %d: requesting floor %d", e.ID, floorNumber)
	e.requestCh <- floorNumber
}

func (e *Elevator) GetUpdates() <-chan struct{} {
	return e.updateCh
}

func (e *Elevator) Lock() {
	e.mu.Lock()
}

func (e *Elevator) Unlock() {
	e.mu.Unlock()
}

func (e *Elevator) GetActivity() Activity {
	e.mu.Lock()
	defer e.mu.Unlock()

	return Activity{
		ID:                e.ID,
		CurrentFloor:      e.CurrentFloor,
		DestinationFloors: append([]int(nil), e.DestinationFloors...),
		Status:            e.Status,
	}
}

func (e *Elevator) requestor() {
	for floor := range e.requestCh {
		e.mu.Lock()
		e.DestinationFloors = append(e.DestinationFloors, floor)
		e.mu.Unlock()
		log.Printf("elevator id %d: added destination floor: %d\n", e.ID, floor)
	}
}

func (e *Elevator) updator() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		e.mu.Lock()
		if len(e.DestinationFloors) == 0 {
			e.Status = IDLE
			e.mu.Unlock()
			continue
		}

		target := e.DestinationFloors[0]
		reached := false

		switch {
		case e.CurrentFloor < target:
			e.CurrentFloor++
			e.Status = UP
		case e.CurrentFloor > target:
			e.CurrentFloor--
			e.Status = DOWN
		default:
			e.Status = IDLE
			e.DestinationFloors = e.DestinationFloors[1:]
			reached = true
		}
		e.mu.Unlock()

		// Non-blocking send: if nobody is listening yet, drop the tick
		select {
		case e.updateCh <- struct{}{}:
		default:
		}

		if reached {
			log.Printf("elevator id %d: arrived at %d", e.ID, target)
		}
	}
}
