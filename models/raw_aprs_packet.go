package models

import (
	"github.com/mholt/binding"
)

type RawAprsPacket struct {
	Data   string `json:"data"`
	IsAX25 bool   `json:"is_ax25"`
}

func (cf *RawAprsPacket) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&cf.Data:   "data",
		&cf.IsAX25: "is_ax25",
	}
}
