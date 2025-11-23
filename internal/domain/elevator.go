package domain

type Elevator struct {
	CurrentFloor      int
	DestinationFloors []int
	Status            Status
}

func NewElevator() *Elevator {
	return &Elevator{DestinationFloors: make([]int, 0)}
}
