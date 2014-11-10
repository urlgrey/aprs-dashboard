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

	// clear out existing list if any
	db.Delete("foo")
	length, err := db.ListLength("foo")
	if 0 != length {
		t.Error("List length should be zero", length)
	}

	// push item onto list
	message := &AprsMessage{}
	err = db.PushHead("foo", message)
	if err != nil {
		t.Error("Error while LPUSHing", err)
	}

	// verify item is on list
	length, err = db.ListLength("foo")
	if 1 != length {
		t.Error("List length should be one", length)
	}

	// delete list
	db.Delete("foo")
}
