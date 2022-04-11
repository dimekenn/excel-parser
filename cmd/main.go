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

// @title           excel-service swagger
// @version         1.0
// @description     This is a sample to parse excel files.

// @contact.name   Dinmukhamed Nurbekov
// @contact.email  nurbekovDR@xplanet.ru

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.stage.xplanet.int/parser/
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
