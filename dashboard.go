package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/urlgrey/aprs-dashboard/db"
	"github.com/urlgrey/aprs-dashboard/handlers"
	"github.com/urlgrey/aprs-dashboard/parser"
	"github.com/zencoder/disque-go/disque"
)

func main() {
	n := negroni.New()

	n.Use(negroni.NewStatic(http.Dir("assets")))
	n.Use(negroni.HandlerFunc(handlers.TokenVerificationMiddleware))

	database := db.NewDatabase()
	defer database.Close()
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()

	router := mux.NewRouter()
	disquePool := createDisquePool()
	handlers.InitializeRouterForMessageHandlers(router, aprsParser, disquePool)
	handlers.InitializeRouterForQueryHandlers(router, database)
	n.UseHandler(router)

	n.Run(":3000")
}

func createDisquePool() (pool *disque.DisquePool) {
	queueServer := strings.TrimLeft(os.Getenv("QUEUE_PORT"), "tcp://")
	if queueServer == "" {
		log.Fatal("QUEUE_PORT environment variable is not set, but is required, exiting")
	}

	hosts := []string{queueServer}  // array of 1 or more Disque servers
	cycle := 1000                   // check connection stats every 1000 Fetch's
	capacity := 10                  // initial capacity of the pool
	maxCapacity := 10               // max capacity that the pool can be resized to
	idleTimeout := 15 * time.Minute // timeout for idle connections
	return disque.NewDisquePool(hosts, cycle, capacity, maxCapacity, idleTimeout)
}
