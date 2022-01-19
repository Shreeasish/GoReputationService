package handler

import (
	"bytes"
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
	var jsonStr = []byte(`{"url":"www.bad.com", "score":"10"}`)
	r, _ := http.NewRequest("POST", "/urlinfo/update", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	lines := `http://www.bad.com,2
http://www.google.com,0`
	s := strings.NewReader(lines)
	scorer, _ := scorer.New(s)

	api := New(scorer)
	api.UpdateDomainHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
