package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/urlgrey/aprs-dashboard/db"
	"github.com/urlgrey/aprs-dashboard/models"
	"github.com/urlgrey/aprs-dashboard/parser"
)

func InitializeRouter(m *martini.ClassicMartini) {
	m.Put("/api/v1/message", binding.Bind(models.RawAprsPacket{}), messageHandler)
}

func messageHandler(message models.RawAprsPacket) (int, []byte) {
	redisHost := strings.TrimLeft(os.Getenv("DB_PORT"), "tcp://")
	if redisHost == "" {
		log.Fatal("DB_PORT environment variable is not set, but is required, exiting")
	}
	redisPassword := os.Getenv("APRS_REDIS_PASSWORD")
	redisDatabase := os.Getenv("APRS_REDIS_DATABASE")

	database := db.NewDatabase(redisHost, redisPassword, redisDatabase)
	defer database.Close()

	aprsParser := parser.NewParser()
	defer aprsParser.Finish()

	aprsMessage, parseErr := aprsParser.ParseAprsPacket(message.Data, message.IsAX25)
	if parseErr != nil {
		body, _ := json.Marshal("{}")
		return http.StatusNotAcceptable, body
	} else {
		if aprsMessage.SourceCallsign != "" {
			database.RecordMessage(aprsMessage.SourceCallsign, aprsMessage)
			body, _ := json.Marshal("{}")
			return http.StatusOK, body
		} else {
			log.Println("Unable to find source callsign in APRS message")
			body, _ := json.Marshal("{}")
			return http.StatusNotAcceptable, body
		}
	}
}
