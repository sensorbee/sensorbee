package config

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"pfi/sensorbee/sensorbee/data"
)

// Config is a root of configuration trees.
type Config struct {
	// Network section has parameters related to the network such as listening
	// ports, timeouts, etc..
	Network *Network

	// Topologies section has information of topologies created on startup.
	// This section might be going to be removed when SensorBee supports
	// persisting topologies.
	Topologies Topologies

	// Storage section has information of storage of components in SensorBee.
	Storage *Storage

	// Logging section has parameters related to logging.
	Logging *Logging
}

var (
	rootSchemaString = fmt.Sprintf(`{
	"type": "object",
	"properties": {
		"network": %v,
		"topologies": %v,
		"storage": %v,
		"logging": %v
	},
	"additionalProperties": false
}`, networkSchemaString, topologiesSchemaString, storageSchemaString, loggingSchemaString)
	rootSchema *gojsonschema.Schema
)

func init() {
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(rootSchemaString))
	if err != nil {
		panic(err)
	}
	rootSchema = s
}

func New(m data.Map) (*Config, error) {
	if err := validate(rootSchema, m); err != nil {
		return nil, err
	}
	return &Config{
		Network:    newNetwork(mustAsMap(getWithDefault(m, "network", data.Map{}))),
		Topologies: newTopologies(mustAsMap(getWithDefault(m, "topologies", data.Map{}))),
		Storage:    newStorage(mustAsMap(getWithDefault(m, "storage", data.Map{}))),
		Logging:    newLogging(mustAsMap(getWithDefault(m, "logging", data.Map{}))),
	}, nil
}

// TODO: Add FromJSON or FromYAML if necessary
