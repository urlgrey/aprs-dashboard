package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-martini/martini"
)

func main() {
	apiTokens := strings.Split(os.Getenv("APRS_API_TOKENS"), ",")
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

	// m.Get("/api/v1/callsign/:callsign", func(req *http.Request, params martini.Params) (int, []byte) {
	// 	page := parseOptionalIntParam(req.URL.Query().Get("page"), 1)
	// 	records, err := db.GetRecordsForCallsign(params["callsign"], page)
	// 	if err == nil {
	// 		body, _ := json.Marshal(records)
	// 		return http.StatusOK, body
	// 	} else {
	// 		log.Println("Unable to find callsign data", params["callsign"])
	// 		body, _ := json.Marshal("{}")
	// 		return http.StatusNotFound, body
	// 	}
	// })
	// m.Get("/api/v1/position", func(req *http.Request, params martini.Params) (int, []byte) {
	// 	var parseErr error
	// 	lat, parseErr := parseRequiredFloatParam(req.URL.Query().Get("lat"))
	// 	long, parseErr := parseRequiredFloatParam(req.URL.Query().Get("long"))
	// 	time := parseOptionalIntParam(req.URL.Query().Get("time"), 3600)
	// 	radiusKM := parseOptionalIntParam(req.URL.Query().Get("radius"), 30)

	// 	if parseErr != nil {
	// 		log.Println("Unable to parse required parameters for lat-long")
	// 		body, _ := json.Marshal("{}")
	// 		return http.StatusBadRequest, body
	// 	} else {
	// 		records, err := db.GetRecordsNearPosition(lat, long, time, radiusKM)

	// 		if err == nil {
	// 			body, _ := json.Marshal(records)
	// 			return http.StatusOK, body
	// 		} else {
	// 			body, _ := json.Marshal("{}")
	// 			return http.StatusInternalServerError, body
	// 		}
	// 	}
	// })
	m.Run()
}
