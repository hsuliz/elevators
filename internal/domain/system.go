package domain

import (
	"log"

	"github.com/hsuliz/elevators/internal/domain/types"
)

type ElevatorActivity struct {
	ID                int
	CurrentFloor      int
	DestinationFloors []int
	Status            types.Status
}

type System struct {
	Elevators  []*Elevator
	Floors     []bool
	Picker     Picker
	ActivityCh chan ElevatorActivity
}

func NewSystem(elevators []*Elevator, floors int, picker Picker) *System {
	system := &System{
		Elevators:  elevators,
		Floors:     make([]bool, floors),
		Picker:     picker,
		ActivityCh: make(chan ElevatorActivity, 100),
	}

	for _, e := range system.Elevators {
		go e.Run()
	}

	system.collectActivity()

	return system
}

func (s *System) Call(floorNumber int) {
	pickedElevatorID := s.Picker.Pick(s.Elevators)
	s.Elevators[pickedElevatorID].requests <- floorNumber
}

func (s *System) GetFloors() []bool {
	return s.Floors
}

func (s *System) collectActivity() {
	for _, e := range s.Elevators {
		e := e
		go func() {
			for range e.updateCh {

				e.mu.Lock()
				snapshot := ElevatorActivity{
					ID:           e.ID,
					CurrentFloor: e.CurrentFloor,
					Status:       e.Status,
				}
				if len(e.DestinationFloors) > 0 {
					snapshot.DestinationFloors = append([]int(nil), e.DestinationFloors...)
				}
				e.mu.Unlock()

				log.Printf(
					"elevator %d updated: floor=%d status=%v",
					snapshot.ID,
					snapshot.CurrentFloor,
					snapshot.Status,
				)

				select {
				case s.ActivityCh <- snapshot:
				default:
				}
			}
		}()
	}
}
