package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/urlgrey/aprs-dashboard/models"
	"github.com/urlgrey/aprs-dashboard/parser"
	"github.com/zencoder/disque-go/disque"
)

type MessageHandler struct {
	parser *parser.AprsParser
	disque disque.Disque
}

func InitializeRouterForMessageHandlers(r *mux.Router, parser *parser.AprsParser) {
	m := MessageHandler{parser: parser}
	m.Initialize()

	r.HandleFunc("/api/v1/message", m.SubmitAPRSMessage).Methods("PUT")
}

func (m *MessageHandler) Initialize() (err error) {
	disqueHostsEnv := os.Getenv("DISQUE_HOSTS")
	var disqueHosts []string
	if disqueHostsEnv == "" {
		disqueHosts = []string{"127.0.0.1:7711"}
	} else {
		disqueHosts = strings.Split(disqueHostsEnv, ",")
	}
	d := disque.NewDisque(disqueHosts, 1000)
	return d.Initialize()
}

func (m *MessageHandler) SubmitAPRSMessage(resp http.ResponseWriter, req *http.Request) {
	message := new(models.RawAprsPacket)
	errs := binding.Bind(req, message)
	if errs.Handle(resp) {
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

	if err = m.disque.Push("aprs_messages", string(aprsMessageJson), 100); err != nil {
		log.Printf("Error while enqueueing APRS message for asynchronous handling: %s", err)
		http.Error(resp,
			"Error queueing APRS message for asynchronous handling",
			http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	responseEncoder := json.NewEncoder(resp)
	responseEncoder.Encode("{}")
}
