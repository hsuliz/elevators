package domain

import (
	"log"
	"sync"
	"testing"
)

func TestElevator_TurnOn(t *testing.T) {
	elevator := NewElevator(0)
	elevator.TurnOn()

	wg := sync.WaitGroup{}
	wg.Go(func() {
		for range elevator.GetUpdates() {
			log.Println(elevator.CurrentFloor, elevator.Status)
		}
	})

	wg.Wait()
}
