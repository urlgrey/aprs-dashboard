package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
)

func TestTokenVerificationMiddlewareGet(t *testing.T) {
	h := NewTokenVerificationMiddleware()
	h.Initialize()

	r, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)
	rw.Header().Set("Content-Type", "application/json")

	h.Run(rw, r, DummyHandler)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTokenVerificationMiddlewarePutWithoutTokenAndNoTokenEnforcement(t *testing.T) {
	h := NewTokenVerificationMiddleware()
	h.Initialize()

	r, _ := http.NewRequest("PUT", "/", strings.NewReader(""))
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)
	rw.Header().Set("Content-Type", "application/json")

	h.Run(rw, r, DummyHandler)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTokenVerificationMiddlewarePutWithoutTokenAndTokenEnforcement(t *testing.T) {
	os.Setenv("APRS_API_TOKENS", "abc123,foobar")
	defer os.Setenv("APRS_API_TOKENS", "")

	h := NewTokenVerificationMiddleware()
	h.Initialize()

	r, _ := http.NewRequest("PUT", "/", strings.NewReader(""))
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)
	rw.Header().Set("Content-Type", "application/json")

	h.Run(rw, r, DummyHandler)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTokenVerificationMiddlewarePutWithTokenAndTokenEnforcement(t *testing.T) {
	os.Setenv("APRS_API_TOKENS", "abc123,foobar")
	defer os.Setenv("APRS_API_TOKENS", "")

	h := NewTokenVerificationMiddleware()
	h.Initialize()

	r, _ := http.NewRequest("PUT", "/", strings.NewReader(""))
	r.Header.Add("Accept", "application/json")
	r.Header.Add("X-API-KEY", "abc123")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)
	rw.Header().Set("Content-Type", "application/json")

	h.Run(rw, r, DummyHandler)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func DummyHandler(resp http.ResponseWriter, req *http.Request) {
}
