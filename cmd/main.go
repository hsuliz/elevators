package main

import (
	"context"
	"errors"
	"github.com/hsuliz/elevators/internal/api"
	"github.com/hsuliz/elevators/internal/domain"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	elevatorCount = 3
	floorCount    = 10
	addr          = ":8080"
)

func main() {
	elevators := make([]*domain.Elevator, elevatorCount)
	for i := range elevators {
		elevators[i] = domain.NewElevator(i + 1)
		elevators[i].TurnOn()
	}

	system := domain.NewSystem(elevators, floorCount)

	srv := api.New(system)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv.WatchAndBroadcast(ctx)

	httpSrv := &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("shutting down…")
		cancel()
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutCancel()
		_ = httpSrv.Shutdown(shutCtx)
	}()

	log.Printf("listening on %s", addr)
	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v", err)
	}
}
