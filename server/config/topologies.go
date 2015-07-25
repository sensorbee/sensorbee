package config

import (
	"github.com/xeipuuv/gojsonschema"
	"pfi/sensorbee/sensorbee/data"
)

// Topology has configuration parameters for a SensorBee topology.
type Topology struct {
	// Name is the name of the topology. This field isn't directly used
	// in a config file.
	Name string `json:"-" yaml:"-"`

	// BQLFile is a file path to the BQL file executed on start up.
	BQLFile string `json:"bql_file" yaml:"bql_file"`
}

// Topologies is a set of configuration of topologies.
type Topologies map[string]*Topology

var (
	topologiesSchemaString = `{
	"type": "object",
	"properties": {
	},
	"patternProperties": {
		".*": {
			"anyOf": [
				{
					"type": "object",
					"properties": {
						"bql_file": {
							"type": "string",
							"minLength": 1
						}
					},
					"additionalProperties": false
				},
				{
					"type": "null"
				}
			]
		}
	}
}`

	// Because gojsonschema doesn't support partial schema validation, this
	// has to be defined separately from
	topologiesSchema *gojsonschema.Schema
)

func init() {
	// TODO: need pattern validation on bql_file if possible
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(topologiesSchemaString))
	if err != nil {
		panic(err)
	}
	topologiesSchema = s
}

// NewTopologies creates a Topologies config parameters from a given map.
func NewTopologies(m data.Map) (Topologies, error) {
	if err := validate(topologiesSchema, m); err != nil {
		return nil, err
	}
	return newTopologies(m), nil
}

func newTopologies(m data.Map) Topologies {
	ts := Topologies{}
	for name, conf := range m {
		if conf.Type() == data.TypeNull {
			conf = data.Map{}
		}
		t := &Topology{
			Name:    name,
			BQLFile: mustAsString(getWithDefault(mustAsMap(conf), "bql_file", data.String(""))),
		}
		ts[name] = t
	}
	return ts
}
