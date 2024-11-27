package main

import (
	"backend/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
)

func main() {
	// routes using the SetupRoutes function from routes.go
	r := routes.SetupRoutes()

	corsOptions := handlers.AllowedOrigins([]string{"*"})

	// Set up server and define shutdown
	server := &http.Server{
		Handler:      handlers.CORS(corsOptions)(r),
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start the server
	go func() {
		log.Println("Server is running on port 8000...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up channel for shutdown on interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-stop

	log.Println("Shutting down the server...")

	// Gracefully shut down the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}
