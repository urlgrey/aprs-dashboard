package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"time"
)

type AprsMessage struct {
	// TODO: change timestamp data type to long
	Timestamp   int64   `json:"timestamp"`
	Callsign    string  `json:"callsign"`
	Status      string  `json:"status"`
	Symbol      string  `json:"symbol"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	Speed       float64 `json:"speed"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	RawMessage  string  `json:"raw_message"`
}

func main() {
	redisHost := os.Getenv("APRS_REDIS_HOST")
	if redisHost == "" {
		log.Fatal("APRS_REDIS_HOST environment variable is not set, but is required, exiting")
	}
	redisPassword := os.Getenv("APRS_REDIS_PASSWORD")
	redisDatabase := os.Getenv("APRS_REDIS_DATABASE")

	redisPool := newPool(redisHost, redisPassword, redisDatabase)
	defer redisPool.Close()

	c := redisPool.Get()
	_, err := c.Do("PING")
	c.Close()
	if err != nil {
		log.Fatal("Error pinging Redis", err)
	}

	m := martini.Classic()
	m.Put("/api/v1/message", binding.Bind(AprsMessage{}), func(message AprsMessage) string {
		return "OK\n"
	})
	m.Run()
}

func newPool(server string, password string, database string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if database != "" {
				if _, err := c.Do("SELECT", database); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
