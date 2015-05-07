package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/urlgrey/aprs-dashboard/models"
	"github.com/urlgrey/aprs-dashboard/parser"
	"github.com/zencoder/disque-go/disque"
	"golang.org/x/net/context"
)

type MessageHandler struct {
	parser      *parser.AprsParser
	pool        *disque.DisquePool
	PoolTimeout time.Duration
}

const (
	APRS_MESSAGES_QUEUE_NAME = "aprs_messages"
)

func NewMessageHandler(parser *parser.AprsParser, pool *disque.DisquePool) *MessageHandler {
	return &MessageHandler{parser: parser, pool: pool, PoolTimeout: 5 * time.Second}
}

func InitializeRouterForMessageHandlers(r *mux.Router, parser *parser.AprsParser, pool *disque.DisquePool) {
	m := NewMessageHandler(parser, pool)
	r.HandleFunc("/api/v1/message", m.SubmitAPRSMessage).Methods("PUT")
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
	ctx := context.Background()
	context.WithTimeout(ctx, m.PoolTimeout)
	if conn, err = m.pool.Get(ctx); err != nil {
		log.Printf("Error while getting connection from pool: %s", err)
		http.Error(resp,
			"Error queueing APRS message for asynchronous handling",
			http.StatusInternalServerError)
		return
	}

	if _, err = conn.Push(APRS_MESSAGES_QUEUE_NAME, string(aprsMessageJson), 100); err != nil {
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
