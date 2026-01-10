package domain

import (
	"log"
	"time"
)

type System struct {
	Elevators  []*Elevator
	Picker     Picker
	Floors     []*Floor
	CallChs    []chan int
	MoveCh     chan int
	ActivityCh chan Elevator
}

func NewSystem(elevators []*Elevator, picker Picker, floorCount int) *System {
	floors := make([]*Floor, floorCount)
	for i := range floorCount {
		floors[i] = NewFloor()
	}

	callChans := make([]chan int, len(elevators))
	for i := range len(elevators) {
		callChans[i] = make(chan int)
	}

	system := &System{
		Elevators:  elevators,
		Picker:     picker,
		Floors:     floors,
		CallChs:    callChans,
		MoveCh:     make(chan int, 100),
		ActivityCh: make(chan Elevator, 100),
	}

	for i := range len(elevators) {
		system.monitor(i)
	}
	system.Activity()

	return system
}

func (s *System) Call(floorNumber int) {
	pickedElevatorId := s.Picker.Pick(s.Elevators)
	s.CallChs[pickedElevatorId] <- floorNumber
}

func (s *System) Activity() {
	go func() {
		for id := range s.MoveCh {
			elevator := s.Elevators[id]
			log.Println("elevator:", elevator, "updated")
			s.ActivityCh <- *elevator
		}
	}()
}

func (s *System) monitor(elevatorId int) {
	go func() {
		elevator := s.Elevators[elevatorId]
		for destinationFloor := range s.CallChs[elevatorId] {
			log.Print("monitor: ", elevatorId, elevator.CurrentFloor, destinationFloor)
			elevator.mu.Lock()
			elevator.DestinationFloors = append(elevator.DestinationFloors, destinationFloor)
			elevator.mu.Unlock()
			// #TODO DestinationFloors self-balancing??
			if len(elevator.DestinationFloors) == 1 {
				go s.move(elevatorId)
			}
		}
	}()
}

func (s *System) move(elevatorId int) {
	elevator := s.Elevators[elevatorId]

	for len(elevator.DestinationFloors) != 0 {
		switch {
		case elevator.DestinationFloors[0] > elevator.CurrentFloor:
			elevator.mu.Lock()
			elevator.CurrentFloor++
			elevator.Status = UP
			elevator.mu.Unlock()
		case elevator.DestinationFloors[0] < elevator.CurrentFloor:
			elevator.mu.Lock()
			elevator.CurrentFloor--
			elevator.Status = DOWN
			elevator.mu.Unlock()
		case elevator.DestinationFloors[0] == elevator.CurrentFloor:
			elevator.mu.Lock()
			elevator.Status = IDLE
			elevator.mu.Unlock()
			s.MoveCh <- elevatorId
			elevator.DestinationFloors = elevator.DestinationFloors[1:]
			log.Printf(
				"elevatorId: %d, came to destination floor: %d",
				elevatorId,
				elevator.CurrentFloor,
			)
			time.Sleep(time.Second)
			continue
		}

		s.MoveCh <- elevatorId
		time.Sleep(time.Second)
		log.Printf("elevatorId: %d, currentFloor: %d", elevatorId, elevator.CurrentFloor)
	}
	log.Print("movement finished")
}
