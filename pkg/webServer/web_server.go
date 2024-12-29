package webServer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	eventStatuses = make(map[string]string)
	mu            sync.Mutex
)

func AddEventStatus(eventID, status string) {
	mu.Lock()
	defer mu.Unlock()
	eventStatuses[eventID] = status
}

func GetEventStatuses(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventStatuses)
}

func StartWebServer(handler http.Handler) {
	http.HandleFunc("/qgo/events", GetEventStatuses)
	http.Handle("/", http.FileServer(http.Dir("web")))

	port := ":6666"
	fmt.Println("Starting web server on http://localhost" + port)
	// Create a server with a timeout for graceful shutdown
	server := &http.Server{Addr: port}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		// Block until we get a signal
		<-sigCh
		fmt.Println("\nReceived shutdown signal, shutting down the server...")

		// Set a timeout for shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Gracefully shutdown the server
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("Error during shutdown: %v\n", err)
		} else {
			fmt.Println("Server shut down gracefully.")
		}

	}()

	if err := http.ListenAndServe(port, handler); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to start web server: %v\n", err)
	}
}
