package domain

type Picker interface {
	Pick(elevators []*Elevator) int
}
