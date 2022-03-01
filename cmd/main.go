package main

import (
	"context"
	"excel-service/internal/transport/http"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//channel for errors
	errCh := make(chan error, 1)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%v", <-sigCh)
	}()

	go http.StartHTTPServer(ctx, errCh)

	log.Fatalf("service terminated: %v", <-errCh)
}
