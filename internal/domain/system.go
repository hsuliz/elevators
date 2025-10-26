package domain

import (
	"log"
	"time"

	"github.com/hsuliz/elevators/internal/domain/picker"
)

type System struct {
	Elevators []*Elevator
	Picker    picker.Interface
	Floors    []*Floor
}

func NewSystem(elevators []*Elevator, pickerInterface picker.Interface, floorsNumber int) *System {
	var floors = make([]*Floor, floorsNumber)
	for i := range floorsNumber {
		floors[i] = NewFloor()
	}
	return &System{
		Elevators: elevators,
		Picker:    pickerInterface,
		Floors:    floors,
	}
}

func (s *System) Call(floorNumber int) bool {
	log.Printf("trying to call elevator from floor number %d", floorNumber)
	if s.IsCalled(floorNumber) {
		log.Printf("elevator already called from floor number %d", floorNumber)
		return false
	}
	s.Floors[floorNumber].Called = true
	log.Printf("called elevator from floor number %d", floorNumber)
	s.Pick(floorNumber)
	return true
}

func (s *System) Pick(floorNumber int) {
	pickedElevator := s.Picker.Pick(len(s.Elevators))
	log.Printf("elevator picked %d", pickedElevator)
	s.Move(pickedElevator, floorNumber)
}

func (s *System) IsCalled(floorNumber int) bool {
	return s.Floors[floorNumber].Called
}

func (s *System) Move(elevatorId int, destinationFloor int) {
	elevator := s.Elevators[elevatorId]
	for elevator.CurrentFlor != destinationFloor {
		log.Printf("elevator %d current flor %d", elevatorId, elevator.CurrentFlor)

		elevator.CurrentFlor++
		time.Sleep(1 * time.Second)
	}
	log.Printf("elevator %d arrived at destiationd floor %d", elevatorId, destinationFloor)
}
