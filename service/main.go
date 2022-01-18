package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var domains = []string{"www.baddomain.com", "www.baddomain2.com"}

type HealthResponse struct {
	Count       int    `json:"DomainCount"`
	LastUpdated string `json:"LastUpdateTime"`
}

func main() {
	NewRouter([]Handler{Handler{"/health", "GET", handleHealthCheck}})
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
