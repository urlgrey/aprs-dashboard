package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urlgrey/aprs-dashboard/db"
)

type QueryHandler struct {
	database *db.Database
}

func InitializeRouterForQueryHandlers(r *mux.Router, database *db.Database) {
	h := &QueryHandler{database: database}
	r.HandleFunc("/api/v1/callsign/:callsign", h.callsignQueryHandler)
	r.HandleFunc("/api/v1/position", h.positionQueryHandler)
}

func (h *QueryHandler) callsignQueryHandler(resp http.ResponseWriter, req *http.Request) {
	page := parseOptionalIntParam(req.URL.Query().Get("page"), 1)
	var records *db.PaginatedCallsignResults
	var err error
	if records, err = h.database.GetRecordsForCallsign(req.URL.Query().Get("callsign"), page); err != nil {
		http.Error(resp,
			fmt.Sprintf("Unable to find callsign data %s", req.URL.Query().Get("callsign")),
			http.StatusNoContent)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	responseEncoder := json.NewEncoder(resp)
	responseEncoder.Encode(records)
}

func (h *QueryHandler) positionQueryHandler(resp http.ResponseWriter, req *http.Request) {
	var err error
	var lat, long float64
	if lat, err = parseRequiredFloatParam(req.URL.Query().Get("lat")); err != nil {
		http.Error(resp,
			"Error parsing latitude query parameter from request",
			http.StatusBadRequest)
		return
	}
	if long, err = parseRequiredFloatParam(req.URL.Query().Get("long")); err != nil {
		http.Error(resp,
			"Error parsing longitude query parameter from request",
			http.StatusBadRequest)
		return
	}
	time := parseOptionalIntParam(req.URL.Query().Get("time"), 3600)
	radiusKM := parseOptionalIntParam(req.URL.Query().Get("radius"), 30)

	var records *db.PositionResults
	if records, err = h.database.GetRecordsNearPosition(lat, long, time, radiusKM); err != nil {
		http.Error(resp,
			fmt.Sprintf("Error looking up APRS records %+v", err),
			http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	responseEncoder := json.NewEncoder(resp)
	responseEncoder.Encode(records)
}
