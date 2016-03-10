package config

import (
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

// Storage has storage configuration parameters for components in SensorBee.
type Storage struct {
	UDS UDSStorage `json:"uds" yaml:"uds"`
}

// UDSStorage has configuration parameters for the storage of UDSs.
type UDSStorage struct {
	Type   string   `json:"type" yaml:"params"`
	Params data.Map `json:"params" yaml:"params"`
}

// Because data.Map doesn't support YAML encoding, UDSStorage.Params has type
// map[string]interface{} instead of data.Map.

var (
	storageSchemaString = `{
	"type": "object",
	"properties": {
		"uds": {
			"anyOf": [
				{
					"type": "object",
					"properties": {
						"type": {
							"enum": ["in_memory"]
						},
						"params": {
							"anyOf": [
								{
									"type": "object",
									"maxProperties": 0
								},
								{
									"type": "null"
								}
							]
						}
					},
					"required": ["type"],
					"additionalProperties": false
				},
				{
					"type": "object",
					"properties": {
						"type": {
							"enum": ["fs"]
						},
						"params": {
							"anyOf": [
								{
									"type": "object",
									"properties": {
										"dir": {
											"type": "string"
										},
										"temp_dir": {
											"type": "string"
										}
									},
									"required": ["dir"],
									"additionalProperties": false
								},
								{
									"type": "null"
								}
							]
						}
					},
					"required": ["type"],
					"additionalProperties": false
				}
			]
		}
	},
	"additionalProperties": false
}`
	storageSchema *gojsonschema.Schema
	// TODO: add patterns for filepath validation
)

func init() {
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(storageSchemaString))
	if err != nil {
		panic(err)
	}
	storageSchema = s
}

// NewStorage creates a Storage config parameters from a given map.
func NewStorage(m data.Map) (*Storage, error) {
	if err := validate(storageSchema, m); err != nil {
		return nil, err
	}
	return newStorage(m), nil
}

func newStorage(m data.Map) *Storage {
	udsParams := getWithDefault(m, "uds.params", data.Map{})
	if udsParams.Type() == data.TypeNull {
		udsParams = data.Map{}
	}

	// Some parameter validation such as a test for existence of a directory
	// should be done in each UDSStorage.

	return &Storage{
		UDS: UDSStorage{
			Type:   mustAsString(getWithDefault(m, "uds.type", data.String("in_memory"))),
			Params: mustAsMap(udsParams),
		},
	}
}

// ToMap returns storage config information as data.Map.
func (s *Storage) ToMap() data.Map {
	return data.Map{
		"uds": data.Map{
			"params": s.UDS.Params,
			"type":   data.String(s.UDS.Type),
		},
	}

}
