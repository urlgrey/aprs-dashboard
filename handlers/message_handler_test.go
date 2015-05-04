package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
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
	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
}

func TestHandleMessageWithDisqueUnreachable(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()
	os.Setenv("QUEUE_PORT", "tcp://127.0.0.1:8811")
	h := MessageHandler{parser: aprsParser}
	assert.Nil(t, h.Initialize())

	r, _ := http.NewRequest("PUT", "/api/v1/message", strings.NewReader("{\"data\":\"K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5\", \"is_ax25\":false}"))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)

	h.SubmitAPRSMessage(rw, r)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandleMessageWithUnparsableAPRSPacket(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()
	os.Setenv("QUEUE_PORT", "tcp://127.0.0.1:7711")
	h := MessageHandler{parser: aprsParser}
	assert.Nil(t, h.Initialize())

	r, _ := http.NewRequest("PUT", "/api/v1/message", strings.NewReader("{\"data\":\"foo\", \"is_ax25\":false}"))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)

	h.SubmitAPRSMessage(rw, r)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleMessage(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()
	os.Setenv("QUEUE_PORT", "tcp://127.0.0.1:7711")
	h := MessageHandler{parser: aprsParser}
	assert.Nil(t, h.Initialize())

	r, _ := http.NewRequest("PUT", "/api/v1/message", strings.NewReader("{\"data\":\"K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5\", \"is_ax25\":false}"))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	rw := negroni.NewResponseWriter(rec)

	h.SubmitAPRSMessage(rw, r)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func BenchmarkHandleMessage(b *testing.B) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()
	os.Setenv("QUEUE_PORT", "tcp://127.0.0.1:7711")
	h := MessageHandler{parser: aprsParser}
	h.Initialize()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, _ := http.NewRequest("PUT", "/api/v1/message", strings.NewReader("{\"data\":\"K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5\", \"is_ax25\":false}"))
			r.Header.Add("Content-Type", "application/json")
			r.Header.Add("Accept", "application/json")
			rec := httptest.NewRecorder()
			rw := negroni.NewResponseWriter(rec)

			h.SubmitAPRSMessage(rw, r)
			assert.Equal(b, http.StatusOK, rec.Code)
		}
	})
}
