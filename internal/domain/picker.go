package domain

import "github.com/hsuliz/elevators/internal/domain/dto"

type Picker interface {
	Pick(elevators []*dto.Elevator) int
}
