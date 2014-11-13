package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	apiTokens := strings.Split(os.Getenv("APRS_API_TOKENS"), ",")

	db := NewDatabase(redisHost, redisPassword, redisDatabase)
	defer db.Close()

	parser := NewParser()
	defer parser.Finish()

	m := martini.Classic()

	m.Use(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "PUT" {
			suppliedApiToken := req.Header.Get("X-API-KEY")
			found := false
			for _, token := range apiTokens {
				if suppliedApiToken == token {
					found = true
					break
				}
			}
			if !found {
				res.WriteHeader(http.StatusUnauthorized)
			}
		}
	})

	m.Put("/api/v1/message", binding.Bind(RawAprsPacket{}), func(message RawAprsPacket) (int, []byte) {
		aprsMessage, parseErr := parser.parseAprsPacket(message.Data, message.IsAX25)
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
	m.Get("/api/v1/callsign/:callsign", func(req *http.Request, params martini.Params) (int, []byte) {
		var page int64 = 1
		pageStr := req.URL.Query().Get("page")
		if pageStr != "" {
			var parseErr error
			page, parseErr = strconv.ParseInt(pageStr, 10, 64)
			if parseErr != nil {
				page = 1
			}
		}

		records, err := db.GetRecordsForCallsign(params["callsign"], page)
		if err == nil {
			body, _ := json.Marshal(records)
			return http.StatusOK, body
		} else {
			log.Println("Unable to find callsign data", params["callsign"])
			body, _ := json.Marshal("{}")
			return http.StatusNotFound, body
		}
	})
	m.Run()
}
