package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/urlgrey/aprs-dashboard/models"
	"github.com/urlgrey/aprs-dashboard/parser"
	"github.com/zencoder/disque-go/disque"
)

type MessageHandler struct {
	parser *parser.AprsParser
	pool   *disque.DisquePool
}

func InitializeRouterForMessageHandlers(r *mux.Router, parser *parser.AprsParser) {
	m := MessageHandler{parser: parser}
	m.Initialize()

	r.HandleFunc("/api/v1/message", m.SubmitAPRSMessage).Methods("PUT")
}

func (m *MessageHandler) Initialize() (err error) {
	queueServer := strings.TrimLeft(os.Getenv("QUEUE_PORT"), "tcp://")
	if queueServer == "" {
		log.Fatal("QUEUE_PORT environment variable is not set, but is required, exiting")
	}

	hosts := []string{queueServer}  // array of 1 or more Disque servers
	cycle := 1000                   // check connection stats every 1000 Fetch's
	capacity := 10                  // initial capacity of the pool
	maxCapacity := 10               // max capacity that the pool can be resized to
	idleTimeout := 15 * time.Minute // timeout for idle connections
	m.pool = disque.NewDisquePool(hosts, cycle, capacity, maxCapacity, idleTimeout)

	return nil
}

func (m *MessageHandler) SubmitAPRSMessage(resp http.ResponseWriter, req *http.Request) {
	message := new(models.RawAprsPacket)
	errs := binding.Bind(req, message)
	if errs.Handle(resp) {
		log.Printf("Error while binding request to model: %s", errs.Error())
		return
	}

	var aprsMessage *models.AprsMessage
	var err error
	if aprsMessage, err = m.parser.ParseAprsPacket(message.Data, message.IsAX25); err != nil {
		http.Error(resp,
			fmt.Sprintf("Error parsing APRS message %+v", err),
			http.StatusBadRequest)
		return
	}

	var aprsMessageJson []byte
	if aprsMessageJson, err = json.Marshal(aprsMessage); err != nil {
		log.Printf("Error serializing parsed APRS message for queueing: %s", err)
		http.Error(resp,
			"Error serializing APRS message for queueing",
			http.StatusInternalServerError)
		return
	}

	var conn *disque.Disque
	if conn, err = m.pool.Get(); err != nil {
		log.Printf("Error while getting connection from pool: %s", err)
		http.Error(resp,
			"Error queueing APRS message for asynchronous handling",
			http.StatusInternalServerError)
		return
	}

	if err = conn.Push("aprs_messages", string(aprsMessageJson), 100); err != nil {
		m.pool.Put(conn)
		log.Printf("Error while enqueueing APRS message for asynchronous handling: %s", err)
		http.Error(resp,
			"Error queueing APRS message for asynchronous handling",
			http.StatusInternalServerError)
		return
	}
	m.pool.Put(conn)

	resp.Header().Set("Content-Type", "application/json")
	responseEncoder := json.NewEncoder(resp)
	responseEncoder.Encode("{}")
}
