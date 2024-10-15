package main

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/kanat_9999/homework/notifier/internal/config"
	"gitlab.ozon.dev/kanat_9999/homework/notifier/internal/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	go handleSignals(cancel)

	notifierService := service.NewNotifier(cfg)
	wg.Add(1)
	go notifierService.Start(ctx, cfg, wg)

	wg.Wait()
}

func handleSignals(cancel context.CancelFunc) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stopChan)

	select {
	case sig := <-stopChan:
		fmt.Printf("received signal: %v, shutting down", sig)
		cancel()
	}
}
