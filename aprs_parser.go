package main

// #cgo pkg-config: --libs libfap
// #include <fap.h>
import "C"

type AprsMessage struct {
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

func parseAprsPacket(message string, isAX25 bool) *AprsMessage {
	return &AprsMessage{}
}
