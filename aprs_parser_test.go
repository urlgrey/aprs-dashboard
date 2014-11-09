package main

import (
	"testing"
)

func Test_ParseAprsPacketWithLocation(t *testing.T) {
	msg := parseAprsPacket("K2YNT>APRS,N2YGK-9*,TRACE3-3:!4033.24N/07421.21W", true)
	if msg.SourceCallsign != "K2YNT" {
		t.Error("Source callsign not parsed correctly")
	}
	if msg.DestinationCallsign != "APRS" {
		t.Error("Destination callsign not parsed correctly")
	}
}

func Benchmark_ParseAprsPacket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseAprsPacket("K2YNT>APRS,N2YGK-9*,TRACE3-3:!4033.24N/07421.21W", true)
	}
}
