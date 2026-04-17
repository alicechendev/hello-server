package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var (
	requestCount int
	mu           sync.Mutex
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	requestCount++
	count := requestCount
	mu.Unlock()

	fmt.Fprintf(w, "Hello, World! (request #%d)\n", count)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count := requestCount
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":   "ok",
		"requests": count,
	})
}

// statsLogger periodically logs request counts in a background goroutine.
func statsLogger(done <-chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			count := requestCount
			mu.Unlock()
			log.Printf("stats: total_requests=%d", count)
		case <-done:
			return
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	done := make(chan struct{})
	go statsLogger(done)

	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		close(done)
		log.Fatal(err)
	}
}
