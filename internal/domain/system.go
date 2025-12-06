package domain

import (
	"log"
	"time"
)

type System struct {
	Elevators []*Elevator
	Picker    Picker
	Floors    []*Floor
	CallChans []chan int
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
		Elevators: elevators,
		Picker:    picker,
		Floors:    floors,
		CallChans: callChans,
	}

	for i := range len(elevators) {
		go system.monitor(i)
	}

	return system
}

func (s *System) Call(floorNumber int) {
	pickedElevatorId := s.Picker.Pick(s.Elevators)
	s.CallChans[pickedElevatorId] <- floorNumber
}

func (s *System) Status() <-chan Elevator {
	ch := make(chan Elevator)

	for _, elevator := range s.Elevators {
		go func() {
			prevState := *elevator
			for {
				current := *elevator

				if !current.Equal(prevState) {
					ch <- current
					prevState = current
				}

				// polling delay
				time.Sleep(50 * time.Millisecond)
			}
		}()
	}

	return ch
}

func (s *System) monitor(elevatorId int) {
	elevator := s.Elevators[elevatorId]
	for floorNumber := range s.CallChans[elevatorId] {
		log.Print("monitor ", elevatorId, floorNumber)
		elevator.DestinationFloors = append(elevator.DestinationFloors, floorNumber)
		// #TODO DestinationFloors self-balancing??
		if len(elevator.DestinationFloors) == 1 {
			go s.move(elevatorId)
		}
	}
}

func (s *System) move(elevatorId int) {
	elevator := s.Elevators[elevatorId]

	for len(elevator.DestinationFloors) != 0 {
		switch {
		case elevator.DestinationFloors[0] > elevator.CurrentFloor:
			elevator.CurrentFloor++
			elevator.Status = UP
		case elevator.DestinationFloors[0] < elevator.CurrentFloor:
			elevator.CurrentFloor--
			elevator.Status = DOWN
		case elevator.DestinationFloors[0] == elevator.CurrentFloor:
			elevator.Status = IDLE
			elevator.DestinationFloors = elevator.DestinationFloors[1:]
			log.Printf(
				"elevatorId: %d, came to destination floor: %d",
				elevatorId,
				elevator.CurrentFloor,
			)
			continue
		}
		log.Printf("elevatorId: %d, currentFloor: %d", elevatorId, elevator.CurrentFloor)

		// simulate activity
		time.Sleep(time.Second)
	}
	log.Print("movement finished")
}
