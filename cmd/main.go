package main

import (
	"os"
	"sync"
	"time"

	"github.com/hsuliz/elevators/internal/domain"
)

func main() {
	elevator1 := domain.NewElevator()
	elevator2 := domain.NewElevator()
	elevators := []*domain.Elevator{elevator1, elevator2}

	naivePicker := domain.NewNaive()

	system := domain.NewSystem(elevators, naivePicker, 11)

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		system.Call(5)
	})
	time.Sleep(time.Second)
	wg.Go(func() {
		system.Call(3)
	})
	wg.Go(func() {
		system.Call(7)
	})
	wg.Wait()
	time.Sleep(time.Second * 5)
	os.Exit(0)
}
