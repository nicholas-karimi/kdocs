package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicholas-karimi/kdocs/internal/web"
)

func main() {

	router := web.NewRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		fmt.Println("KDocs server starting on  http://localhost:8080")
		serverErr <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		fmt.Println(err)

	case <-shutdownSignal:
		fmt.Println("Shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Graceful shutdown failed:", err)
	}

}
