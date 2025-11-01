package dto

type Floor struct {
	Called bool
}

func NewFloor() *Floor {
	return &Floor{}
}
