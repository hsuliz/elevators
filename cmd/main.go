package main

import (
	"log"
	"sync"

	"github.com/hsuliz/elevators/internal/domain"
	"github.com/hsuliz/elevators/internal/domain/dto"
)

func main() {
	elevator1 := dto.NewElevator()
	elevator2 := dto.NewElevator()
	elevators := []*dto.Elevator{elevator1, elevator2}

	rngPicker := domain.NewNaive()

	system := domain.NewSystem(elevators, rngPicker, 11)

	//systemElevators := system.Elevators
	//for _, e := range systemElevators {
	//	log.Print(e)
	//}

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		system.Call(3)
	})

	wg.Go(func() {
		system.Call(5)
	})
	wg.Wait()

	log.Print(system.Status())
	wg.Go(func() {
		system.Call(10)
	})
	wg.Wait()
	log.Print(system.Status())
}
