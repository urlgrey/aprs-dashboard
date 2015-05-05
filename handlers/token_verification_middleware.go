package handlers

import (
	"net/http"
	"os"
	"strings"
)

type TokenVerificationMiddleware struct {
	apiTokens []string
}

func NewTokenVerificationMiddleware() *TokenVerificationMiddleware {
	return &TokenVerificationMiddleware{}
}

func (t *TokenVerificationMiddleware) Initialize() {
	t.apiTokens = strings.Split(os.Getenv("APRS_API_TOKENS"), ",")
}

func (t *TokenVerificationMiddleware) Run(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "PUT" {
		suppliedApiToken := r.Header.Get("X-API-KEY")
		found := false
		for _, token := range t.apiTokens {
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
