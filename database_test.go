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

func Test_RecordMessage(t *testing.T) {
	db := NewDatabase(":6379", "", "")
	defer db.Close()

	// clear out existing list if any
	db.Delete("foo")
	db.Delete("callsigns.set")

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

	// delete list
	db.Delete("callsign.foo")
	db.Delete("callsigns.set")
}
