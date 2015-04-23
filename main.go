package main

import (
	"encoding/json"
	"errors"
        "fmt"
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
	redisHost := os.Getenv("DB_PORT_6379_TCP_ADDR")
	if redisHost == "" {
		log.Fatal("APRS_REDIS_HOST environment variable is not set, but is required, exiting")
	}
	redisPassword := os.Getenv("APRS_REDIS_PASSWORD")
	redisDatabase := os.Getenv("APRS_REDIS_DATABASE")
	apiTokens := strings.Split(os.Getenv("APRS_API_TOKENS"), ",")

	db := NewDatabase(fmt.Sprintf("%s:6379", redisHost), redisPassword, redisDatabase)
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

	m.Use(martini.Static("assets")) // serve from the "assets" directory as well

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
		page := parseOptionalIntParam(req.URL.Query().Get("page"), 1)
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
	m.Get("/api/v1/position", func(req *http.Request, params martini.Params) (int, []byte) {
		var parseErr error
		lat, parseErr := parseRequiredFloatParam(req.URL.Query().Get("lat"))
		long, parseErr := parseRequiredFloatParam(req.URL.Query().Get("long"))
		time := parseOptionalIntParam(req.URL.Query().Get("time"), 3600)
		radiusKM := parseOptionalIntParam(req.URL.Query().Get("radius"), 30)

		if parseErr != nil {
			log.Println("Unable to parse required parameters for lat-long")
			body, _ := json.Marshal("{}")
			return http.StatusBadRequest, body
		} else {
			records, err := db.GetRecordsNearPosition(lat, long, time, radiusKM)

			if err == nil {
				body, _ := json.Marshal(records)
				return http.StatusOK, body
			} else {
				body, _ := json.Marshal("{}")
				return http.StatusInternalServerError, body
			}
		}
	})
	m.Run()
}

func parseOptionalIntParam(val string, defaultValue int64) int64 {
	valInt, parseErr := strconv.ParseInt(val, 10, 64)
	if parseErr != nil {
		valInt = defaultValue
	}
	return valInt
}

func parseRequiredFloatParam(val string) (float64, error) {
	valFloat, parseErr := strconv.ParseFloat(val, 64)
	if parseErr != nil {
		return valFloat, errors.New("Float parameter was missing or could not be parsed")
	}
	return valFloat, parseErr
}
