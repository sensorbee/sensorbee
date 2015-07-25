package config

import (
	"encoding/json"
	"pfi/sensorbee/sensorbee/data"
)

func toMap(js string) data.Map {
	var m data.Map
	err := json.Unmarshal([]byte(js), &m)
	if err != nil {
		panic(err)
	}
	return m
}
