package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
	"github.com/urlgrey/aprs-dashboard/parser"
)

func TestHandlerInitialize(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()

	m := MessageHandler{parser: aprsParser}
	err := m.Initialize()
	assert.Nil(t, err)
}

func TestHandlerSubmitAPRSMessage(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()

	m := MessageHandler{parser: aprsParser}
	err := m.Initialize()
	assert.Nil(t, err)
}

func TestHandleMessageWithMalformedJSON(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()
	h := MessageHandler{parser: aprsParser}
	assert.Nil(t, h.Initialize())

	r, _ := http.NewRequest("PUT", "/api/v1/message", strings.NewReader("asdf"))
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)
	rw.Header().Set("Content-Type", "application/json")

	h.SubmitAPRSMessage(rw, r)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleMessageWithDisqueUnreachable(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()
	h := MessageHandler{parser: aprsParser}
	assert.Nil(t, h.Initialize())

	r, _ := http.NewRequest("PUT", "/api/v1/message", strings.NewReader("asdf"))
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)

	h.SubmitAPRSMessage(rw, r)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
