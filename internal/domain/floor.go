package domain

type Floor struct {
	Called bool
}

func NewFloor() *Floor {
	return &Floor{}
}
