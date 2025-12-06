package domain

type Elevator struct {
	ID                int
	CurrentFloor      int
	DestinationFloors []int
	Status            Status
}

func NewElevator(id int) *Elevator {
	return &Elevator{ID: id, DestinationFloors: make([]int, 0)}
}

func (e Elevator) Equal(other Elevator) bool {
	return e.ID == other.ID &&
		e.CurrentFloor == other.CurrentFloor &&
		e.Status == other.Status
}
