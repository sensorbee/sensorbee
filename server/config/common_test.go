package config

import (
	"encoding/json"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

func toMap(js string) data.Map {
	var m data.Map
	err := json.Unmarshal([]byte(js), &m)
	if err != nil {
		panic(err)
	}
	return m
}
