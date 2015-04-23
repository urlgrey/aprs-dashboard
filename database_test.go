package main

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func Test_NewDatabase(t *testing.T) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	err := db.Ping()
	if err != nil {
		t.Error("Unable to ping Redis", err)
	}
}

func Test_RecordMessage(t *testing.T) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	var err error
	var length int64

	// verify item is not on list
	length, err = db.NumberOfMessagesForCallsign("KK6DCI")
	if 0 != length {
		t.Error("List length should be one", length)
	}

	// verify item is not on list
	length, err = db.NumberOfCallsigns()
	if 0 != length {
		t.Error("List length should be one", length)
	}

	// push item onto list
	message := &AprsMessage{SourceCallsign: "KK6DCI"}
	err = db.RecordMessage("KK6DCI", message)
	if err != nil {
		t.Error("Error while LPUSHing", err)
	}

	// verify item is stored
	var mostRecentMsg *AprsMessage
	mostRecentMsg, err = db.GetMostRecentMessageForCallsign("KK6DCI")
	if err != nil {
		t.Error("Error while getting most record message for callsign", err)
	}
	if mostRecentMsg.SourceCallsign != "KK6DCI" {
		t.Error("Most recent message for callsign was invalid, missing callsign", mostRecentMsg.SourceCallsign)
	}

	// verify item is on list
	length, err = db.NumberOfMessagesForCallsign("KK6DCI")
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

func Test_GetMostRecentMessageForCallsign(t *testing.T) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("KK6DCI>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)
	db.RecordMessage(msg.SourceCallsign, msg)

	mostRecentMsg, err := db.GetMostRecentMessageForCallsign("KK6DCI")
	if err != nil {
		t.Error("Unexpected error")
	}
	if mostRecentMsg.SourceCallsign != "KK6DCI" {
		t.Error("Most recent message reported incomplete")
	}
}

func Test_GetMostRecentMessageForUnrecognizedCallsign(t *testing.T) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	mostRecentMsg, err := db.GetMostRecentMessageForCallsign("FOADFBASDF")
	if err == nil {
		t.Error("Expected error but saw none")
	}
	if mostRecentMsg.SourceCallsign != "" {
		t.Error("Most recent message reported incomplete")
	}
}

func Test_GetRecordsNearPosition(t *testing.T) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	for i := 1; i <= 10; i++ {
		msg, _ := p.parseAprsPacket("KK6DCI-"+strconv.Itoa(i)+">APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)
		db.RecordMessage(msg.SourceCallsign, msg)
	}

	nearbyRecords, err := db.GetRecordsNearPosition(47.720333333333336, -122.3735, 3600, 30)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if nearbyRecords.Size != 10 {
		t.Error("Size of response reported incorrectly", nearbyRecords.Size)
	}
}

func Benchmark_RecordMessage(b *testing.B) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("KK6DCI>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	for i := 0; i < b.N; i++ {
		db.RecordMessage(msg.SourceCallsign, msg)
	}
}
func Benchmark_RetrieveMostRecentEntriesForCallsign(b *testing.B) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("KK6DCI>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	var i int
	for i = 0; i < 10000; i++ {
		db.RecordMessage("KK6DCI", msg)
	}

	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		db.GetRecordsForCallsign("KK6DCI", 1)
	}
}

func Benchmark_RetrieveMiddleEntriesForCallsign(b *testing.B) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	msg, _ := p.parseAprsPacket("KK6DCI>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)

	var i int
	for i = 0; i < 10000; i++ {
		db.RecordMessage("KK6DCI", msg)
	}

	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		db.GetRecordsForCallsign("KK6DCI", 500)
	}
}

func Benchmark_GetFormattedTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFormattedTime(time.Now())
	}
}

func Benchmark_GetRecordsNearPosition(b *testing.B) {
	db := NewDatabase(os.Getenv("DB_PORT_6379_TCP_ADDR"), "", "")
	defer db.Close()
	cleanup(db)
	defer cleanup(db)

	p := NewParser()
	defer p.Finish()
	for i := 1; i <= 10; i++ {
		msg, _ := p.parseAprsPacket("KK6DCI-"+strconv.Itoa(i)+">APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)
		db.RecordMessage(msg.SourceCallsign, msg)
	}

	for i := 0; i < b.N; i++ {
		db.GetRecordsNearPosition(47.720333333333336, -122.3735, 3600, 30)
	}
}

func cleanup(db *Database) {
	db.Delete("aprs:calls")
	db.Delete("aprs:last:KK6DCI")
	db.Delete("aprs:pos")
	db.Delete("aprs:all:KK6DCI")
	for i := 1; i <= 10; i++ {
		db.Delete("aprs:last:KK6DCI-" + strconv.Itoa(i))
		db.Delete("aprs:all:KK6DCI-" + strconv.Itoa(i))
	}
}
