package domain

import (
	"log"
	"time"

	"github.com/hsuliz/elevators/internal/domain/dto"
)

type System struct {
	Elevators []*dto.Elevator
	Picker    Picker
	Floors    []*dto.Floor
}

func NewSystem(elevators []*dto.Elevator, picker Picker, floorsNumber int) *System {
	var floors = make([]*dto.Floor, floorsNumber)
	for i := range floorsNumber {
		floors[i] = dto.NewFloor()
	}
	return &System{
		Elevators: elevators,
		Picker:    picker,
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
	elevatorsAsValue := make([]dto.Elevator, len(s.Elevators))
	for i, e := range s.Elevators {
		elevatorsAsValue[i] = *e
	}
	pickedElevator := s.Picker.Pick(elevatorsAsValue)
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
		time.Sleep(time.Millisecond)
	}
	log.Printf("elevator %d arrived at destiationd floor %d", elevatorId, destinationFloor)
}

func (s *System) Status() []dto.Elevator {
	elevatorsAsValue := make([]dto.Elevator, len(s.Elevators))
	for i, e := range s.Elevators {
		elevatorsAsValue[i] = *e
	}
	return elevatorsAsValue
}
