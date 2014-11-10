package main

import (
	"testing"
)

func Test_ParseAprsPacketWithLocation(t *testing.T) {
	msg := parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", true)
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

func Benchmark_ParseAprsPacket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseAprsPacket("K7SSW>APRS,TCPXX*,qAX,CWOP-5:@100235z4743.22N/12222.41W_135/000g000t047r004p009P008h95b10132lOww_0.86.5", true)
	}
}
