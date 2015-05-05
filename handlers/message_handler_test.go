package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
	"github.com/urlgrey/aprs-dashboard/parser"
	"github.com/zencoder/disque-go/disque"
)

func TestHandleMessageWithMalformedJSON(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()
	pool := createDisquePool("127.0.0.1:7711")
	defer pool.Close()
	h := NewMessageHandler(aprsParser, pool)

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
	pool := createDisquePool("127.0.0.1:8811")
	defer pool.Close()
	h := NewMessageHandler(aprsParser, pool)

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
	pool := createDisquePool("127.0.0.1:7711")
	defer pool.Close()
	h := NewMessageHandler(aprsParser, pool)

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
	pool := createDisquePool("127.0.0.1:7711")
	defer pool.Close()
	h := NewMessageHandler(aprsParser, pool)

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
	pool := createDisquePool("127.0.0.1:7711")
	defer pool.Close()
	h := NewMessageHandler(aprsParser, pool)

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

func createDisquePool(server string) (pool *disque.DisquePool) {
	hosts := []string{server}       // array of 1 or more Disque servers
	cycle := 1000                   // check connection stats every 1000 Fetch's
	capacity := 10                  // initial capacity of the pool
	maxCapacity := 10               // max capacity that the pool can be resized to
	idleTimeout := 15 * time.Minute // timeout for idle connections
	return disque.NewDisquePool(hosts, cycle, capacity, maxCapacity, idleTimeout)
}
