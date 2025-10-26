package main

import (
	"sync"

	"github.com/hsuliz/elevators/internal/domain"
	"github.com/hsuliz/elevators/internal/domain/picker"
)

func main() {
	elevator1 := domain.NewElevator()
	elevator2 := domain.NewElevator()
	elevator3 := domain.NewElevator()
	elevator4 := domain.NewElevator()
	elevator5 := domain.NewElevator()
	elevators := []*domain.Elevator{elevator1, elevator2, elevator3, elevator4, elevator5}

	rngPicker := picker.NewRNG()

	system := domain.NewSystem(elevators, rngPicker, 10)

	//systemElevators := system.Elevators
	//for _, e := range systemElevators {
	//	log.Print(e)
	//}

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		system.Call(3)
	})
	wg.Go(func() {
		system.Call(2)
	})
	wg.Wait()
}
