package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicholas-karimi/kdocs/internal/config"
	"github.com/nicholas-karimi/kdocs/internal/web"
)

func main() {

	// read config
	cfg := config.Load()
	router, err := web.NewRouter()

	if err != nil {
		fmt.Println("Failed to initialize router:", err)
		return
	}
	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		fmt.Printf("KDocs server starting on  http://localhost%s\n", cfg.Addr)
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
