package config

import (
	"github.com/xeipuuv/gojsonschema"
	"pfi/sensorbee/sensorbee/data"
)

// Network has configuration parameters related to the network.
type Network struct {
	// ListenOn has binding information in "host:port" format.
	ListenOn string `json:"listen_on" yaml:"listen_on"`
}

var networkSchema *gojsonschema.Schema

func init() {
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(`
{
	"type": "object",
	"properties": {
		"listen_on": {
			"type": "string",
			"pattern": "^.*:[0-9]+$"
		}
	},
	"additionalProperties": false
}
`))
	if err != nil {
		panic(err)
	}
	networkSchema = s
}

// NewNetwork creates a Newtork config parameters from a given map.
func NewNetwork(m data.Map) (*Network, error) {
	if err := validate(networkSchema, m); err != nil {
		return nil, err
	}

	return &Network{
		ListenOn: mustAsString(getWithDefault(m, "listen_on", data.String(":8090"))),
	}, nil
}
