package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shreeasish/reputation/scorer"
	"github.com/gorilla/mux"
)

// ApiHandler is a struct blah blah
type ApiHandler struct {
	rp *scorer.DomainScorer
}

func New(r *scorer.DomainScorer) *ApiHandler {
	return &ApiHandler{
		rp: r,
	}
}

func sendBadRequestResponse(w http.ResponseWriter, m string) {
	type BadRequestResponse struct {
		Message string `json:"Message"`
	}

	w.WriteHeader(http.StatusBadRequest)
	b := BadRequestResponse{
		Message: m,
	}
	res, _ := json.Marshal(b)
	w.Write(res)
}

// GetScore to handle calls to retrieve information about a Domain
func (h ApiHandler) GetScoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url, ok := vars["url"]
	if !ok {
		sendBadRequestResponse(w, "Required parameter 'url' not found")
		return
	}

	ds := h.rp.GetDomainScore(url)

	res, err := json.Marshal(ds)
	if err != nil {
		log.Printf("Unable to marshal json. Err %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

type UpdateDomainHandlerResponse struct {
	Message string `json:"Message"`
}

type Update struct {
	Url   string `json:"url"`
	Score string `json:"score"`
}

// UpdateDomainHandler handles requests to update or add new domains with scores
func (h ApiHandler) UpdateDomainHandler(w http.ResponseWriter, r *http.Request) {
	var u Update
	if r.Body == nil {
		log.Printf("Request with invalid body")
		sendBadRequestResponse(w, "Invalid request body")
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Printf("Unable to parse message body %v", err)
		sendBadRequestResponse(w, "Unable to parse request payload")
		return
	}
	defer r.Body.Close()

	i, err := strconv.ParseInt(u.Score, 10, 64)
	if err != nil {
		log.Printf("Unable to parse score %v", err)
		sendBadRequestResponse(w, fmt.Sprintf("Invalid domain score %v", err))
		return
	}
	h.rp.AddDomain(u.Url, i)

	res, err := json.Marshal(UpdateDomainHandlerResponse{
		Message: "Domain updated successfully",
	})
	if err != nil {
		log.Printf("Unable to marshal json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

type HealthResponse struct {
	Count       int       `json:"DomainCount"`
	LastUpdated time.Time `json:"LastUpdateTime"`
}

func (h ApiHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	hr := HealthResponse{
		Count:       h.rp.Count,
		LastUpdated: h.rp.LastUpdated,
	}
	res, err := json.Marshal(hr)
	if err != nil {
		log.Printf("Unable to marshal json. Err %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
