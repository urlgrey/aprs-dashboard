package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"log"
	"os"
)

func main() {
	redisHost := os.Getenv("APRS_REDIS_HOST")
	if redisHost == "" {
		log.Fatal("APRS_REDIS_HOST environment variable is not set, but is required, exiting")
	}
	redisPassword := os.Getenv("APRS_REDIS_PASSWORD")
	redisDatabase := os.Getenv("APRS_REDIS_DATABASE")

	db := NewDatabase(redisHost, redisPassword, redisDatabase)
	defer db.Close()

	m := martini.Classic()
	m.Put("/api/v1/message", binding.Bind(AprsMessage{}), func(message AprsMessage) string {
		return "OK\n"
	})
	m.Run()
}
