package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urlgrey/aprs-dashboard/parser"
)

func TestHandlerInitialize(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()

	m := MessageHandler{parser: aprsParser}
	err := m.Initialize()
	assert.Nil(t, err)
}

func TestHandlerSubmitAPRSMessage(t *testing.T) {
	aprsParser := parser.NewParser()
	aprsParser.Initialize()
	defer aprsParser.Close()

	m := MessageHandler{parser: aprsParser}
	err := m.Initialize()
	assert.Nil(t, err)
}
