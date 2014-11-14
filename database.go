package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"time"
)

type Database struct {
	redisPool *redis.Pool
}

type PositionResults struct {
	Size    int64         `json:"size"`
	Records []AprsMessage `json:"records"`
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
	}

	return &db
}

func (db *Database) Close() {
	db.redisPool.Close()
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

func (db *Database) RecordMessage(sourceCallsign string, message *AprsMessage) error {
	jsonBytes, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		log.Println("Error converting message to JSON", marshalErr)
		return nil
	} else {
		c := db.redisPool.Get()
		defer c.Close()

		msgString := string(jsonBytes[:])
		numberOfCommands := 3
		c.Send("HINCRBY", "callsigns.set", sourceCallsign, 1)
		c.Send("LPUSH", "callsign."+sourceCallsign, msgString)
		c.Send("SET", "callsign.lastmessage."+sourceCallsign, msgString)
		if message.IncludesPosition {
			numberOfCommands = numberOfCommands + 2
			c.Send("geoadd", "positions", message.Latitude, message.Longitude, sourceCallsign)
			c.Send("geoadd", "positions."+getFormattedTime(time.Now()), message.Latitude, message.Longitude, sourceCallsign)
		}
		c.Flush()

		var err error
		for i := 0; i < numberOfCommands; i++ {
			_, err = c.Receive()
		}

		return err
	}
}

func (db *Database) GetMostRecentMessageForCallsign(callsign string) (*AprsMessage, error) {
	c := db.redisPool.Get()
	defer c.Close()

	msgBytes, err := redis.Bytes(c.Do("GET", "callsign.lastmessage."+callsign))
	if err == nil {
		var m AprsMessage
		err = json.Unmarshal(msgBytes, &m)
		return &m, err
	} else {
		return &AprsMessage{}, err
	}
}

func (db *Database) NumberOfMessagesForCallsign(callsign string) (int64, error) {
	c := db.redisPool.Get()
	defer c.Close()
	r, err := c.Do("LLEN", "callsign."+callsign)
	return r.(int64), err
}

func (db *Database) PaginatedMessagesForCallsign(callsign string, start int64, stop int64) ([]string, error) {
	c := db.redisPool.Get()
	defer c.Close()
	return redis.Strings(redis.Values(c.Do("LRANGE", "callsign."+callsign, start, stop)))
}

func (db *Database) NumberOfCallsigns() (int64, error) {
	c := db.redisPool.Get()
	defer c.Close()
	r, err := c.Do("HLEN", "callsigns.set")
	return r.(int64), err
}

func (db *Database) GetRecordsNearPosition(lat float64, long float64, timeInterval int64, radiusKM int64) (*PositionResults, error) {
	c := db.redisPool.Get()
	defer c.Close()

	currentSearchTime := time.Now().Truncate(time.Duration(1) * time.Hour)
	numberOfSearches := int(timeInterval / 3600)
	for i := 0; i < numberOfSearches; i++ {
		redis.Strings(redis.Values(c.Do("georadius", "positions."+getFormattedTime(currentSearchTime), lat, long, strconv.FormatInt(radiusKM, 10)+" km")))
		currentSearchTime = currentSearchTime.Add(time.Duration(1) * time.Hour)
	}

	return &PositionResults{}, nil
}

func (db *Database) GetRecordsForCallsign(callsign string, page int64) (*PaginatedCallsignResults, error) {
	var err error
	totalNumberOfRecords, err := db.NumberOfMessagesForCallsign(callsign)
	if err == nil {
		numberOfPages := (totalNumberOfRecords / 10) + 1
		startingRecord := (page - 1) * 10
		endingRecord := (page * 10) - 1
		resultingMessages := []AprsMessage{}
		messages, _ := db.PaginatedMessagesForCallsign(callsign, startingRecord, endingRecord)
		for _, message := range messages {
			var m AprsMessage
			unmarshalErr := json.Unmarshal([]byte(message), &m)
			if unmarshalErr == nil {
				resultingMessages = append(resultingMessages, m)
			} else {
				log.Println("Unable to parse message, skipping")
			}
		}
		results := PaginatedCallsignResults{
			Page:                 page,
			NumberOfPages:        numberOfPages,
			TotalNumberOfRecords: totalNumberOfRecords,
			Records:              resultingMessages,
		}
		return &results, nil
	} else {
		return &PaginatedCallsignResults{}, errors.New("Unable to get the number of records for the specified callsign")
	}
}

func getFormattedTime(t time.Time) string {
	utcTime := t.UTC()
	return fmt.Sprintf("%d.%02d.%02d.%02d",
		utcTime.Year(), utcTime.Month(), utcTime.Day(),
		utcTime.Hour())
}
