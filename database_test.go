package main

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func Test_NewDatabase(t *testing.T) {
	db := NewDatabase(os.Getenv("APRS_REDIS_HOST"), "", "")
	defer db.Close()
	err := db.Ping()
	if err != nil {
		t.Error("Unable to ping Redis", err)
	}
}

func Test_RecordMessage(t *testing.T) {
	db := NewDatabase(os.Getenv("APRS_REDIS_HOST"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	var err error
	var length int64

	// verify item is not on list
	length, err = db.NumberOfMessagesForCallsign("foo")
	if 0 != length {
		t.Error("List length should be one", length)
	}

	// verify item is not on list
	length, err = db.NumberOfCallsigns()
	if 0 != length {
		t.Error("List length should be one", length)
	}

	// push item onto list
	message := &AprsMessage{}
	err = db.RecordMessage("foo", message)
	if err != nil {
		t.Error("Error while LPUSHing", err)
	}

	// verify item is on list
	length, err = db.NumberOfMessagesForCallsign("foo")
	if 1 != length {
		t.Error("List length should be one", length)
	}

	// verify item is on list
	length, err = db.NumberOfCallsigns()
	if 1 != length {
		t.Error("List length should be one", length)
	}
}

func Test_getFormattedTime(t *testing.T) {
	i, err := strconv.ParseInt("1405544146", 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	tm := time.Unix(i, 0)
	timeStr := getFormattedTime(tm)
	if timeStr != "2014.07.16.20" {
		t.Error("Formatted time string is incorrect:", timeStr)
	}
}

func Benchmark_RecordMessage(b *testing.B) {
	db := NewDatabase(os.Getenv("APRS_REDIS_HOST"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	for i := 0; i < b.N; i++ {
		db.RecordMessage(msg.SourceCallsign, msg)
	}
}
func Benchmark_RetrieveMostRecentEntriesForCallsign(b *testing.B) {
	db := NewDatabase(os.Getenv("APRS_REDIS_HOST"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	var i int
	for i = 0; i < 10000; i++ {
		db.RecordMessage("foo", msg)
	}

	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		db.GetRecordsForCallsign("foo", 1)
	}
}

func Benchmark_RetrieveMiddleEntriesForCallsign(b *testing.B) {
	db := NewDatabase(os.Getenv("APRS_REDIS_HOST"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	var i int
	for i = 0; i < 10000; i++ {
		db.RecordMessage("foo", msg)
	}

	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		db.GetRecordsForCallsign("foo", 500)
	}
}

func Benchmark_GetFormattedTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFormattedTime(time.Now())
	}
}

func cleanup(db *Database) {
	db.Delete("callsign.foo")
	db.Delete("callsigns.set")
	db.Delete("positions")
	db.Delete("positions." + getFormattedTime(time.Now()))
}
