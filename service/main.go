package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var domains = []string{"www.baddomain.com", "www.baddomain2.com"}

type HealthResponse struct {
	Count       int    `json:"DomainCount"`
	LastUpdated string `json:"LastUpdateTime"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", handleHealthCheck).Methods("GET")
	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Starting server")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(server)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	hr := HealthResponse{
		Count:       len(domains),
		LastUpdated: time.Now().String(),
	}
	res, err := json.Marshal(hr)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func waitForShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
