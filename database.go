package main

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
	"time"
)

type Database struct {
	redisPool *redis.Pool
}

type PaginatedCallsignResults struct {
	Page                 int64         `json:"page"`
	NumberOfPages        int64         `json:"number_of_pages"`
	TotalNumberOfRecords int64         `json:"total_number_of_records"`
	Records              []AprsMessage `json:"records"`
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

func (db *Database) RecordMessage(key string, message *AprsMessage) error {
	jsonBytes, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		log.Println("Error converting message to JSON", marshalErr)
		return nil
	} else {
		c := db.redisPool.Get()
		defer c.Close()

		var err error
		_, err = c.Do("HINCRBY", "callsigns.set", key, 1)
		if err == nil {
			recordKey := strings.Join([]string{"callsign", key}, ".")
			_, err = c.Do("LPUSH", recordKey, string(jsonBytes[:]))
		}
		return err
	}
}

func (db *Database) Ping() error {
	c := db.redisPool.Get()
	defer c.Close()
	_, err := c.Do("PING")
	return err
}

func (db *Database) Delete(key string) {
	c := db.redisPool.Get()
	defer c.Close()
	c.Do("DEL", key)
}

func (db *Database) NumberOfMessagesForCallsign(key string) (int64, error) {
	c := db.redisPool.Get()
	defer c.Close()
	r, err := c.Do("LLEN", "callsign."+key)
	return r.(int64), err
}

func (db *Database) NumberOfCallsigns() (int64, error) {
	c := db.redisPool.Get()
	defer c.Close()
	r, err := c.Do("HLEN", "callsigns.set")
	return r.(int64), err
}

func (db *Database) GetRecordsForCallsign(callsign string, page string) (PaginatedCallsignResults, error) {
	var err error
	totalNumberOfRecords, err := db.NumberOfMessagesForCallsign(callsign)
	if err == nil {
		numberOfPages := (totalNumberOfRecords / 10) + 1
		records := []AprsMessage{}
		results := PaginatedCallsignResults{
			Page:                 1,
			NumberOfPages:        numberOfPages,
			TotalNumberOfRecords: totalNumberOfRecords,
			Records:              records,
		}
		return results, nil
	} else {
		return PaginatedCallsignResults{}, errors.New("Unable to get the number of records for the specified callsign")
	}
}
