package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/urlgrey/aprs-dashboard/db"
	"github.com/urlgrey/aprs-dashboard/models"
	"github.com/urlgrey/aprs-dashboard/parser"
)

type MessageHandler struct {
	database *db.Database
	parser   *parser.AprsParser
}

func InitializeRouterForMessageHandlers(r *mux.Router, database *db.Database, parser *parser.AprsParser) {
	h := MessageHandler{database: database, parser: parser}
	r.HandleFunc("/api/v1/message", h.messageHandler).Methods("PUT")
}

func (h *MessageHandler) messageHandler(resp http.ResponseWriter, req *http.Request) {
	message := new(models.RawAprsPacket)
	errs := binding.Bind(req, message)
	if errs.Handle(resp) {
		return
	}

	var aprsMessage *models.AprsMessage
	var err error
	if aprsMessage, err = h.parser.ParseAprsPacket(message.Data, message.IsAX25); err != nil {
		http.Error(resp,
			fmt.Sprintf("Error parsing APRS message %+v", err),
			http.StatusInternalServerError)
		return
	}

	if aprsMessage.SourceCallsign == "" {
		http.Error(resp,
			fmt.Sprintf("Unable to find source callsign in APRS message"),
			http.StatusBadRequest)
		return
	}

	if err = h.database.RecordMessage(aprsMessage.SourceCallsign, aprsMessage); err != nil {
		log.Printf("Error while storing APRS message: %s", err)
		http.Error(resp,
			"Error storing APRS message",
			http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	responseEncoder := json.NewEncoder(resp)
	responseEncoder.Encode("{}")
}
