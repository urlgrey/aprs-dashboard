package handlers

import (
	"errors"
	"strconv"
)

func parseOptionalIntParam(val string, defaultValue int64) int64 {
	valInt, parseErr := strconv.ParseInt(val, 10, 64)
	if parseErr != nil {
		valInt = defaultValue
	}
	return valInt
}

func parseRequiredFloatParam(val string) (float64, error) {
	valFloat, parseErr := strconv.ParseFloat(val, 64)
	if parseErr != nil {
		return valFloat, errors.New("Float parameter was missing or could not be parsed")
	}
	return valFloat, parseErr
}
