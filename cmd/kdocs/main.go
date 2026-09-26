package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicholas-karimi/kdocs/internal/config"
	"github.com/nicholas-karimi/kdocs/internal/database"
	"github.com/nicholas-karimi/kdocs/internal/web"
)

func main() {

	// read config
	cfg := config.Load()

	// db
	db, err := database.Open(cfg.DBDSN)
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	router, err := web.NewRouter(db)

	if err != nil {
		log.Println("Failed to initialize router:", err)
		return
	}
	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("KDocs server starting on  http://localhost%s\n", cfg.Addr)
		serverErr <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Println(err)

	case <-shutdownSignal:
		log.Println("Shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Graceful shutdown failed:", err)
	}

}
