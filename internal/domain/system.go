package domain

import (
	"log"
	"time"

	"github.com/hsuliz/elevators/internal/domain/picker"
	"github.com/hsuliz/elevators/internal/domain/types"
)

type System struct {
	Elevators  []*types.Elevator
	Picker     picker.Interface
	Floors     []*types.Floor
	CallChs    []chan int
	MoveCh     chan int
	ActivityCh chan types.Elevator
}

func NewSystem(elevators []*types.Elevator, picker picker.Interface, floorCount int) *System {
	callChans := make([]chan int, len(elevators))
	for i := range len(elevators) {
		callChans[i] = make(chan int)
	}

	system := &System{
		Elevators:  elevators,
		Picker:     picker,
		Floors:     createFloors(floorCount),
		CallChs:    callChans,
		MoveCh:     make(chan int, 100),
		ActivityCh: make(chan types.Elevator, 100),
	}

	for i := range len(elevators) {
		system.monitor(i)
	}
	system.Activity()

	return system
}

func createFloors(floorCount int) []*types.Floor {
	floors := make([]*types.Floor, floorCount)
	for i := range floorCount {
		floors[i] = types.NewFloor()
	}
	return floors
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
			elevator.Mu.Lock()
			elevator.DestinationFloors = append(elevator.DestinationFloors, destinationFloor)
			elevator.Mu.Unlock()
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
			elevator.Mu.Lock()
			elevator.CurrentFloor++
			elevator.Status = types.UP
			elevator.Mu.Unlock()
		case elevator.DestinationFloors[0] < elevator.CurrentFloor:
			elevator.Mu.Lock()
			elevator.CurrentFloor--
			elevator.Status = types.DOWN
			elevator.Mu.Unlock()
		case elevator.DestinationFloors[0] == elevator.CurrentFloor:
			elevator.Mu.Lock()
			elevator.Status = types.IDLE
			elevator.Mu.Unlock()
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
