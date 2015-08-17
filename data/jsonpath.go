package data

import (
	"fmt"
	"strconv"
)

// Path is an entity that can be evaluated with a Map to
// return the value stored at the location specified by
// the Path.
type Path interface {
	evaluate(Map) (Value, error)
	set(Map, Value) error
}

// MustCompilePath takes a JSON Path as a string and returns
// an instance of Path representing that JSON Path, or panics
// if the parameter is not a valid JSON Path.
func MustCompilePath(s string) Path {
	p, err := CompilePath(s)
	if err != nil {
		panic(err.Error())
	}
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
	// parse the statement
	j := &jsonPeg{}
	j.Buffer = s
	j.Init()
	if err := j.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing '%s' as a JSON Path", s)
	}
	j.Execute()
	return j, nil
}

// evaluate returns the entry of the map located at the JSON Path
// represented by this jsonPeg instance.
func (j *jsonPeg) evaluate(m Map) (Value, error) {
	// `current` holds the Value into which we descend, the extracted
	// value is then written to `next` by `c.extract()`. By assigning
	// `current = next` after `c.extract()` returns, we can go deeper.
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

// set sets the entry of the map located at the JSON Path represented
// by this jsonPeg instance. Missing intermediat children will be created
// when needed (`mkdir -p` behavior), but if, say, there is an Int located
// at `foo.bar`, then assigning to `foo.bar.hoge` will fail.
func (j *jsonPeg) set(m Map, v Value) error {
	if m == nil || m.Type() == TypeNull {
		return fmt.Errorf("given Map is inaccessible")
	}
	// `current` holds the Value into which we descend, the extracted
	// value is then written to `next` by `c.extractForSet()`. By assigning
	// `current = next` after `c.extractForSet()` returns, we can go deeper.
	var current Value = m
	var next Value
	// setValueInParent is a closure that
	// - is set in `extractForSet`
	// - when called, writes the given value at the position where `next`
	//   was located when `extractForSet` returned.
	var setValueInParent func(v Value)

	for _, c := range j.components {
		err := c.extractForSet(current, &next, &setValueInParent)
		if err != nil {
			return err
		}
		current = next
	}
	if setValueInParent == nil {
		return fmt.Errorf("setValueInParent was nil, when it shouldn't be")
	}
	setValueInParent(v)
	return nil
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
	// due to parser configuration, s will always be a numeric string,
	// but it may overflow int32, so we need a check here.
	if err != nil {
		// TODO panic is not the gold standard of error handling, but
		//      at the moment we have no better way to signal an error
		//      from within jsonPeg.Execute()
		panic(fmt.Sprintf("overflow index number: " + s))
	}
	j.components = append(j.components, &arrayElementExtractor{int(i)})
}

// extractor describes an entity that can extract a child element
// from a Value.
type extractor interface {
	extract(v Value, next *Value) error
	extractForSet(Value, *Value, *func(Value)) error
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

func (a *mapValueExtractor) extractForSet(v Value, next *Value, setInParent *func(Value)) error {
	// if there is a NULL value where a map is supposed to be,
	// then we will create a map
	if v.Type() == TypeNull {
		v = Map{a.key: Null{}}
		(*setInParent)(v)
	}
	// access as a Map
	cont, err := v.asMap()
	if err != nil {
		return fmt.Errorf("cannot access a %T using key \"%s\"", v, a.key)
	}
	// if the Map does not have the key, add it (so that we
	// can "descend" further into cont[a.key])
	if _, ok := cont[a.key]; !ok {
		cont[a.key] = Null{}
	}
	// invariant: cont[a.key] is a valid entry here, possibly NULL
	*setInParent = func(v Value) {
		cont[a.key] = v
	}
	*next = cont[a.key]
	return nil
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

func (a *arrayElementExtractor) extractForSet(v Value, next *Value, setInParent *func(Value)) error {
	// if there is a NULL value where an array is supposed to be,
	// then we will create an array that holds enough entries
	if v.Type() == TypeNull {
		x := make(Array, a.idx+1)
		for i := range x {
			x[i] = Null{}
		}
		v = x
		(*setInParent)(v)
	}
	// access as an Array
	cont, err := v.asArray()
	if err != nil {
		return fmt.Errorf("cannot access a %T using index %d", v, a.idx)
	}
	// if the Array is not long enough, pad it with NULLs (so
	// that we can "descend" further into cont[a.idx])
	if a.idx >= len(cont) {
		for i := len(cont); i <= a.idx; i++ {
			cont = append(cont, Null{})
		}
		// we need to write the possibly reallocated slice
		// to the correct position
		(*setInParent)(cont)
	}
	// invariant: cont[a.idx] is a valid entry, possibly NULL
	*setInParent = func(v Value) {
		cont[a.idx] = v
	}
	*next = cont[a.idx]
	return nil
}
