package main

import (
	"testing"
)

func Test_parseOptionalIntParam_Missing(t *testing.T) {
	if 6 != parseOptionalIntParam("", 6) {
		t.Error("Parsed optional int incorrectly")
	}
}

func Test_parseOptionalIntParam_Supplied(t *testing.T) {
	if 5 != parseOptionalIntParam("5", 6) {
		t.Error("Parsed optional int incorrectly")
	}
}

func Test_parseRequiredFloatParam_Missing(t *testing.T) {
	_, err := parseRequiredFloatParam("")
	if err == nil {
		t.Error("Parsed required Float incorrectly")
	}
}

func Test_parseRequiredFloatParam_Supplied(t *testing.T) {
	val, err := parseRequiredFloatParam("5.0")
	if val != 5.0 {
		t.Error("Parsed float value was incorrect", val)
	}
	if err != nil {
		t.Error("Unexpected error parsing required float that's present")
	}
}

func Benchmark_parseOptionalIntParam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseOptionalIntParam("5", 10)
	}
}
