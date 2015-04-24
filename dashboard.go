package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-martini/martini"
	"github.com/urlgrey/aprs-dashboard/handlers"
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
	handlers.InitializeRouterForMessageHandlers(m)
	handlers.InitializeRouterForQueryHandlers(m)

	m.Run()
}
