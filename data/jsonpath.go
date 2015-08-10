package data

import (
	"fmt"
	"strconv"
)

// Path is an entity that can be evaluated with a Map to
// return the value stored at the location specified by
// the Path.
type Path interface {
	Evaluate(Map) (Value, error)
}

// MustCompilePath takes a JSON Path as a string and returns
// an instance of Path representing that JSON Path, or panics
// if the parameter is not a valid JSON Path.
func MustCompilePath(s string) Path {
	// parse the statement
	p := &jsonPeg{}
	p.Buffer = s
	p.Init()
	if err := p.Parse(); err != nil {
		panic(err.Error())
	}
	p.Execute()
	return p
}

// CompilePath takes a JSON Path as a string and returns an
// instance of Path representing that JSON Path, or an error
// if the parameter is not a valid JSON Path.
func CompilePath(s string) (p Path, err error) {
	// catch any parser errors
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return MustCompilePath(s), nil
}

// Evaluate returns the entry of the map located at the JSON Path
// represented by this jsonPeg instance.
func (j *jsonPeg) Evaluate(m Map) (Value, error) {
	// now get the contents of the map
	var current Value = m
	var next Value
	for _, c := range j.components {
		err := c.extract(current, &next)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// addMapAccess is called when we discover `foo` or `['bar']`
// in a JSON Path string.
func (j *jsonPeg) addMapAccess(s string) {
	j.components = append(j.components, &mapValueExtractor{s})
}

// addArrayAccess is called when we discover `[1]` in a JSON Path
// string.
func (j *jsonPeg) addArrayAccess(s string) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("overflow index number: " + s))
	}
	j.components = append(j.components, &arrayElementExtractor{int(i)})
}

// extractor describes an entity that can extract a child element
// from a Value.
type extractor interface {
	extract(v Value, next *Value) error
}

// mapValueExtractor can extract a value from a Map using the
// given key.
type mapValueExtractor struct {
	key string
}

func (a *mapValueExtractor) extract(v Value, next *Value) error {
	cont, err := AsMap(v)
	if err != nil {
		return err
	}
	if elem, ok := cont[a.key]; ok {
		*next = elem
		return nil
	}
	return fmt.Errorf("key '%s' was not found in map", a.key)
}

// arrayElementExtractor can extract an element from an Array using
// the given index.
type arrayElementExtractor struct {
	idx int
}

func (a *arrayElementExtractor) extract(v Value, next *Value) error {
	cont, err := AsArray(v)
	if err != nil {
		return fmt.Errorf("cannot access a %T using index %d", v, a.idx)
	}
	if a.idx < len(cont) {
		*next = cont[a.idx]
		return nil
	}
	return fmt.Errorf("out of range access: %d", a.idx)
}
