package domain

import (
	"log"
	"sync"
	"testing"
	"time"
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
	time.Sleep(time.Second)
	elevator.RequestFloor(5)

	//time.Sleep(time.Second * 3)
	//elevator.RequestFloor(3)

	time.Sleep(time.Second * 2)
	elevator.RequestFloor(10)

	wg.Wait()
}
