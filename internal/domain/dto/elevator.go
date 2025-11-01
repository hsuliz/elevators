package dto

type Elevator struct {
	CurrentFlor int
	Status      Status
}

func NewElevator() *Elevator {
	return &Elevator{}
}
