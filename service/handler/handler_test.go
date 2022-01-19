package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Shreeasish/reputation/scorer"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	r, _ := http.NewRequest("GET", "/urlinfo/url/www.google.com", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"url": "www.google.com",
	}

	r = mux.SetURLVars(r, vars)

	lines := `http://www.bad.com,2
http://www.google.com,0`
	s := strings.NewReader(lines)
	scorer, _ := scorer.New(s)

	api := New(scorer)
	api.GetScoreHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateHandler(t *testing.T) {
	r, _ := http.NewRequest("PUT", "/urlinfo/update/score/2/url/www.google.com", nil)
	w := httptest.NewRecorder()

	lines := `http://www.bad.com,2
http://www.google.com,0`
	s := strings.NewReader(lines)
	scorer, _ := scorer.New(s)

	vars := map[string]string{
		"url":   "www.bad.com",
		"score": "2",
	}
	r = mux.SetURLVars(r, vars)

	api := New(scorer)
	api.UpdateDomainHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
