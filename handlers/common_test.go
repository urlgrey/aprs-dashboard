package handlers

import (
	"testing"
)

func TestParseOptionalIntParamMissing(t *testing.T) {
	if 6 != parseOptionalIntParam("", 6) {
		t.Error("Parsed optional int incorrectly")
	}
}

func TestParseOptionalIntParamSupplied(t *testing.T) {
	if 5 != parseOptionalIntParam("5", 6) {
		t.Error("Parsed optional int incorrectly")
	}
}

func TestParseRequiredFloatParam_Missing(t *testing.T) {
	_, err := parseRequiredFloatParam("")
	if err == nil {
		t.Error("Parsed required Float incorrectly")
	}
}

func TestParseRequiredFloatParamSupplied(t *testing.T) {
	val, err := parseRequiredFloatParam("5.0")
	if val != 5.0 {
		t.Error("Parsed float value was incorrect", val)
	}
	if err != nil {
		t.Error("Unexpected error parsing required float that's present")
	}
}

func BenchmarkParseOptionalIntParam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseOptionalIntParam("5", 10)
	}
}
