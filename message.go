package main

import (
	"encoding/json"
)

type message struct{
	Type string `json:"type"`
	User string `json:"user"`
	Text string `json:"text"`
	Color string `json:"color"`
}

func (msg *message) encode() []byte {
	bytes, err := json.Marshal(msg)

	// unlikely
	if err != nil { panic(err) }

	return bytes
}

