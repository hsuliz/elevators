package domain

import (
	"context"
	"fmt"
	"log"
)

type System struct {
	Elevators  []*Elevator
	FloorCount int
	picker     *SmartPicker
}

func NewSystem(elevators []*Elevator, floorCount int) *System {
	return &System{
		Elevators:  elevators,
		FloorCount: floorCount,
		picker:     NewSmartPicker(),
	}
}

func (s *System) Call(floorNumber int) error {
	if len(s.Elevators) == 0 {
		return fmt.Errorf("no elevators available in the system")
	}
	if floorNumber < 0 || floorNumber > s.FloorCount {
		return fmt.Errorf("floor %d out of range [0, %d]", floorNumber, s.FloorCount)
	}

	idx := s.picker.Pick(s.Elevators, floorNumber)
	s.Elevators[idx].RequestFloor(floorNumber)
	return nil
}

func (s *System) MonitorElevator(elevatorID int, ctx context.Context) {
	if elevatorID >= len(s.Elevators) {
		log.Printf("invalid elevator ID: %d", elevatorID)
		return
	}

	elevator := s.Elevators[elevatorID]
	for {
		select {
		case <-elevator.GetUpdates():
			log.Printf("elevator %d current floor: %d", elevator.ID, elevator.CurrentFloor)
		case <-ctx.Done():
			log.Printf("stopping monitor for elevator %d", elevator.ID)
			return
		}
	}
}
