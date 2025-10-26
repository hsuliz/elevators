package picker

import (
	"math/rand/v2"
)

type RNG struct{}

func NewRNG() *RNG { return &RNG{} }

func (r RNG) Pick(elevatorsLen int) int {
	return rand.IntN(elevatorsLen)
}
