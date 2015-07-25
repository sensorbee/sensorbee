package config

import (
	"github.com/xeipuuv/gojsonschema"
	"pfi/sensorbee/sensorbee/data"
)

// Topology has configuration parameters for a SensorBee topology.
type Topology struct {
	// Name is the name of the topology. This field isn't directly used
	// in a config file.
	Name string `json:"name" yaml:"name"`

	// BQLFile is a file path to the BQL file executed on start up.
	BQLFile string `json:"bql_file" yaml:"bql_file"`
}

// Topologies is a set of configuration of topologies.
type Topologies map[string]*Topology

var topologiesSchema *gojsonschema.Schema

func init() {
	// TODO: need pattern validation on bql_file if possible
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(`
{
	"type": "object",
	"properties": {
	},
	"patternProperties": {
		".*": {
			"type": "object",
			"properties": {
				"bql_file": {
					"type": "string",
					"minLength": 1,
				}
			},
			"additionalProperties": false
		}
	}
}`))
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

	ts := Topologies{}
	for name, conf := range m {
		t := &Topology{
			Name:    name,
			BQLFile: mustAsString(getWithDefault(mustAsMap(conf), "bql_file", data.String(""))),
		}
		ts[name] = t
	}
	return ts, nil
}
