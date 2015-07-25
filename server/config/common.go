package config

import (
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"pfi/sensorbee/sensorbee/data"
	"strings"
)

func mustGet(m data.Map, path string) data.Value {
	v, err := m.Get(path)
	if err != nil {
		panic(err)
	}
	return v
}

func getWithDefault(m data.Map, path string, def data.Value) data.Value {
	v, err := m.Get(path)
	if err != nil {
		return def
	}
	return v
}

func mustAsString(v data.Value) string {
	s, err := data.AsString(v)
	if err != nil {
		panic(err)
	}
	return s
}

func mustToBool(v data.Value) bool {
	b, err := data.ToBool(v)
	if err != nil {
		panic(err)
	}
	return b
}

func validate(schema *gojsonschema.Schema, m data.Map) error {
	// GoLoader marshal and unmarshal the map.
	res, err := schema.Validate(gojsonschema.NewGoLoader(m))
	if err != nil {
		return err
	}
	if !res.Valid() {
		// TODO: provide better format
		var errs []string
		for _, e := range res.Errors() {
			errs = append(errs, fmt.Sprintf("- %s", e))
		}
		return errors.New("validation errors:\n" + strings.Join(errs, "\n"))
	}
	return nil
}
