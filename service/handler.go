package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiHandler struct {
	rp *ReputationProvider
}

func NewApiHandler(path string) *ApiHandler {
	return &ApiHandler{
		rp: InitializeFromPath(path),
	}
}

func (h ApiHandler) GetDomainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]
	ds, found := h.rp.GetDomainScore(url)

	// Handle uknown domains with 0
	if !found {
		ds = &DomainScore{
			Domain: url,
			Score:  "0",
		}
	}
	res, err := json.Marshal(ds)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h ApiHandler) UpdateDomainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s := vars["score"]
	url := vars["url"]

	h.rp.AddDomain(&DomainScore{
		Domain: url,
		Score:  s,
	})
}
