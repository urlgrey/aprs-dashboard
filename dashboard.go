package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/awslabs/aws-sdk-go/service/dynamodb"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/urlgrey/aprs-dashboard/db"
	"github.com/urlgrey/aprs-dashboard/handlers"
	"github.com/urlgrey/aprs-dashboard/ingest"
	"github.com/urlgrey/aprs-dashboard/parser"
	"github.com/zencoder/disque-go/disque"
)

func main() {
	mainServer := negroni.New()

	mainServer.Use(negroni.NewStatic(http.Dir("assets")))
	tokenVerification := handlers.NewTokenVerificationMiddleware()
	tokenVerification.Initialize()
	mainServer.Use(negroni.HandlerFunc(tokenVerification.Run))

	database := db.NewDatabase()
	defer database.Close()
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()

	router := mux.NewRouter()
	disquePool := createDisquePool()
	handlers.InitializeRouterForMessageHandlers(router, aprsParser, disquePool)
	handlers.InitializeRouterForQueryHandlers(router, database)
	mainServer.UseHandler(router)

	go mainServer.Run(":3000")

	ingestProcessor := ingest.NewIngestProcessor(disquePool, "aprs_messages", createDynamoDBConnection())
	go ingestProcessor.Run()

	healthCheckServer := negroni.New()
	healthCheckRouter := mux.NewRouter()
	handlers.InitializeRouterForHealthCheckHandler(healthCheckRouter, disquePool)
	healthCheckServer.UseHandler(healthCheckRouter)
	healthCheckServer.Run(":3100")
}

func createDynamoDBConnection() *dynamodb.DynamoDB {
	return dynamodb.New(nil)
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
