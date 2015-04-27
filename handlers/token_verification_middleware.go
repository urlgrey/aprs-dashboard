package handlers

import (
	"net/http"
	"os"
	"strings"
)

func TokenVerificationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	apiTokens := strings.Split(os.Getenv("APRS_API_TOKENS"), ",")
	if r.Method == "PUT" {
		suppliedApiToken := r.Header.Get("X-API-KEY")
		found := false
		for _, token := range apiTokens {
			if suppliedApiToken == token {
				found = true
				break
			}
		}
		if !found {
			rw.WriteHeader(http.StatusUnauthorized)
		}
	}

	next(rw, r)
}
