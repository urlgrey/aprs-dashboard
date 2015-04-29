package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/urlgrey/aprs-dashboard/db"
	"github.com/urlgrey/aprs-dashboard/handlers"
	"github.com/urlgrey/aprs-dashboard/parser"
)

func main() {
	n := negroni.New()

	n.Use(negroni.NewStatic(http.Dir("assets")))
	n.Use(negroni.HandlerFunc(handlers.TokenVerificationMiddleware))

	database := db.NewDatabase()
	defer database.Close()
	aprsParser := parser.NewParser()
	defer aprsParser.Finish()

	router := mux.NewRouter()
	handlers.InitializeRouterForMessageHandlers(router, aprsParser)
	handlers.InitializeRouterForQueryHandlers(router, database)
	n.UseHandler(router)

	n.Run(":3000")
}
