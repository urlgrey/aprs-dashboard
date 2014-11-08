package main

import (
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"log"
	"time"
)

type Database struct {
	redisPool *redis.Pool
}

func NewDatabase(server string, password string, database string) *Database {
	db := Database{}
	db.redisPool = &redis.Pool{
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

	// test the connection now
	c := db.redisPool.Get()
	_, err := c.Do("PING")
	c.Close()
	if err != nil {
		log.Fatal("Error pinging Redis", err)
	} else {
		log.Println("Connection ping succeeded")
	}

	return &db
}

func (db *Database) Close() {
	db.redisPool.Close()
}

func (db *Database) PushHead(key string, message *AprsMessage) error {
	jsonBytes, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		log.Println("Error converting message to JSON", marshalErr)
		return nil
	} else {
		c := db.redisPool.Get()
		defer c.Close()

		_, err := c.Do("LPUSH", key, string(jsonBytes[:]))
		return err
	}
}

func (db *Database) Ping() error {
	c := db.redisPool.Get()
	defer c.Close()
	_, err := c.Do("PING")
	return err
}
