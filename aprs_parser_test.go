package main

import (
	"testing"
	"time"
)

func Test_ParseAprsNonAX25PacketWithLocation(t *testing.T) {
	timeBeforeTest := int32(time.Now().Unix())
	p := NewParser()
	defer p.Finish()

	msg, err := p.parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)
	timeAfterTest := int32(time.Now().Unix())

	if err != nil {
		t.Error("Error was unexpectedly non-nil", err)
	}
	if msg.Timestamp < timeBeforeTest || msg.Timestamp > timeAfterTest {
		t.Error("Timestamp was set incorrectly")
	}
	if msg.SourceCallsign != "K7SSW" {
		t.Error("Source callsign not parsed correctly", msg.SourceCallsign)
	}
	if msg.DestinationCallsign != "APRS" {
		t.Error("Destination callsign not parsed correctly", msg.DestinationCallsign)
	}
	if msg.Latitude != 47.720333333333336 {
		t.Error("Latitude incorrect", msg.Latitude)
	}
	if msg.Longitude != -122.3735 {
		t.Error("Longitude incorrect", msg.Longitude)
	}
}

func Test_ParseAprsNonAX25PacketWithLocationAndWeather(t *testing.T) {
	timeBeforeTest := int32(time.Now().Unix())
	p := NewParser()
	defer p.Finish()

	msg, err := p.parseAprsPacket("DW6161>APRS,TCPXX*,qAX,CWOP-4:@101509z5004.48N/00645.00E_049/000g000t046r000p019P013h97b10123WeatherCatV123B16H31", false)
	timeAfterTest := int32(time.Now().Unix())

	if err != nil {
		t.Error("Error was unexpectedly non-nil", err)
	}
	if msg.Timestamp < timeBeforeTest || msg.Timestamp > timeAfterTest {
		t.Error("Timestamp was set incorrectly")
	}
	if msg.SourceCallsign != "DW6161" {
		t.Error("Source callsign not parsed correctly", msg.SourceCallsign)
	}
	if msg.DestinationCallsign != "APRS" {
		t.Error("Destination callsign not parsed correctly", msg.DestinationCallsign)
	}
	if msg.Latitude != 50.074666666666666 {
		t.Error("Latitude incorrect", msg.Latitude)
	}
	if msg.Longitude != 6.75 {
		t.Error("Longitude incorrect", msg.Longitude)
	}
	if msg.Weather.Temperature != 7.777777777777778 {
		t.Error("Temperature incorrect", msg.Weather.Temperature)
	}
	if msg.Weather.Humidity != 97 {
		t.Error("Humidity incorrect", msg.Weather.Humidity)
	}
	if msg.Weather.InsideHumidity != 0 {
		t.Error("Inside Humidity incorrect", msg.Weather.InsideHumidity)
	}
}

func Test_ParseAprsUnsupportedFormatPacket(t *testing.T) {
	p := NewParser()
	defer p.Finish()

	msg, err := p.parseAprsPacket("ZS6EY>APN982,ZS0TRG*,WIDE3-2,qAS,ZS6EY-1:g {UIV32N}", false)
	if err == nil {
		t.Error("Error was unexpectedly nil", err)
	}
	if msg.SourceCallsign != "" {
		t.Error("Source callsign not parsed correctly", msg.SourceCallsign)
	}
	if msg.DestinationCallsign != "" {
		t.Error("Destination callsign not parsed correctly", msg.DestinationCallsign)
	}
	if msg.Latitude != 0 {
		t.Error("Latitude incorrect", msg.Latitude)
	}
	if msg.Longitude != 0 {
		t.Error("Longitude incorrect", msg.Longitude)
	}
}

func Benchmark_ParseValidAprsPacket(b *testing.B) {
	p := NewParser()
	defer p.Finish()

	for i := 0; i < b.N; i++ {
		p.parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", false)
	}
}

func Benchmark_ParseInvalidAprsPacket(b *testing.B) {
	p := NewParser()
	defer p.Finish()

	for i := 0; i < b.N; i++ {
		p.parseAprsPacket("ZS6EY>APN982,ZS0TRG*,WIDE3-2,qAS,ZS6EY-1:g {UIV32N}", false)
	}
}
