package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	pool := createDisquePool("127.0.0.1:7711")
	defer pool.Close()
	h := NewHealthCheckHandler(pool)

	r, _ := http.NewRequest("GET", "/healthcheck", strings.NewReader(""))
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)
	rw.Header().Set("Content-Type", "application/json")

	h.HealthCheck(rw, r)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthCheckWithBrokenDisquePool(t *testing.T) {
	pool := createDisquePool("127.0.0.1:8811")
	defer pool.Close()
	h := NewHealthCheckHandler(pool)

	r, _ := http.NewRequest("GET", "/healthcheck", strings.NewReader(""))
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)
	rw.Header().Set("Content-Type", "application/json")

	h.HealthCheck(rw, r)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
