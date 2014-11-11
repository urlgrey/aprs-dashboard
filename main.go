package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"log"
	"net/http"
	"os"
)

type RawAprsPacket struct {
	Data   string `json:"data"`
	IsAX25 bool   `json:"is_ax25"`
}

func main() {
	redisHost := os.Getenv("APRS_REDIS_HOST")
	if redisHost == "" {
		log.Fatal("APRS_REDIS_HOST environment variable is not set, but is required, exiting")
	}
	redisPassword := os.Getenv("APRS_REDIS_PASSWORD")
	redisDatabase := os.Getenv("APRS_REDIS_DATABASE")

	db := NewDatabase(redisHost, redisPassword, redisDatabase)
	defer db.Close()

	m := martini.Classic()
	m.Put("/api/v1/message", binding.Bind(RawAprsPacket{}), func(message RawAprsPacket) (int, []byte) {
		aprsMessage, parseErr := parseAprsPacket(message.Data, message.IsAX25)
		if parseErr != nil {
			body, _ := json.Marshal("{}")
			return http.StatusNotAcceptable, body
		} else {
			if aprsMessage.SourceCallsign != "" {
				db.RecordMessage(aprsMessage.SourceCallsign, aprsMessage)
				body, _ := json.Marshal("{}")
				return http.StatusOK, body
			} else {
				log.Println("Unable to find source callsign in APRS message")
				body, _ := json.Marshal("{}")
				return http.StatusNotAcceptable, body
			}
		}
	})
	m.Get("/api/v1/callsign/:callsign", func(params martini.Params) (int, []byte) {
		records, err := db.GetRecordsForCallsign(params["callsign"], "1")
		if err == nil {
			body, _ := json.Marshal(records)
			return http.StatusOK, body
		} else {
			log.Println("Unable to find callsign data")
			body, _ := json.Marshal("{}")
			return http.StatusNotFound, body
		}
	})
	m.Get("/api/v1/callsign/:callsign/:page", func(params martini.Params) (int, []byte) {
		records, err := db.GetRecordsForCallsign(params["callsign"], params["page"])
		if err == nil {
			body, _ := json.Marshal(records)
			return http.StatusOK, body
		} else {
			log.Println("Unable to find callsign data")
			body, _ := json.Marshal("{}")
			return http.StatusNotFound, body
		}
	})
	m.Run()
}
