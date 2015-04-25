package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/urlgrey/aprs-dashboard/db"
)

func InitializeRouterForQueryHandlers(m *martini.ClassicMartini) {
	m.Get("/api/v1/callsign/:callsign", callsignQueryHandler)
	m.Get("/api/v1/position", positionQueryHandler)
}

func callsignQueryHandler(req *http.Request, params martini.Params) (int, []byte) {
	database := db.NewDatabase()
	defer database.Close()

	page := parseOptionalIntParam(req.URL.Query().Get("page"), 1)
	records, err := database.GetRecordsForCallsign(params["callsign"], page)
	if err == nil {
		body, _ := json.Marshal(records)
		return http.StatusOK, body
	} else {
		log.Println("Unable to find callsign data", params["callsign"])
		body, _ := json.Marshal("{}")
		return http.StatusNotFound, body
	}
}

func positionQueryHandler(req *http.Request, params martini.Params) (int, []byte) {
	database := db.NewDatabase()
	defer database.Close()

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
		records, err := database.GetRecordsNearPosition(lat, long, time, radiusKM)

		if err == nil {
			body, _ := json.Marshal(records)
			return http.StatusOK, body
		} else {
			body, _ := json.Marshal("{}")
			return http.StatusInternalServerError, body
		}
	}
}
