package main

import (
	"testing"
)

func Test_NewDatabase(t *testing.T) {
	db := NewDatabase(":6379", "", "")
	defer db.Close()
	err := db.Ping()
	if err != nil {
		t.Error("Unable to ping Redis", err)
	}
}

func Test_PushHead(t *testing.T) {
	db := NewDatabase(":6379", "", "")
	defer db.Close()

	message := &AprsMessage{}
	err := db.PushHead("foo", message)
	if err != nil {
		t.Error("Error while LPUSHing", err)
	}
}
