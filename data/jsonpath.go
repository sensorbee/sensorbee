package data

import (
	"fmt"
	"strconv"
	"strings"
)

type multiplicity int

const (
	one multiplicity = iota
	many
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
	// discover nested array slice accesses
	containsSlice := false
	for _, c := range j.components {
		if c.resultMultiplicity() == many {
			if containsSlice {
				return nil, fmt.Errorf("path '%s' contains multiple slice elements", s)
			}
			containsSlice = true
		}
	}
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

	// resultIsArray is set to true after we processed an
	// extractor that returns an array-valued result by nature (such
	// as arraySliceExtractor). The consequence of `resultIsArray == true`
	// is that we will not process `current` itself with subsequent
	// extractors, but each item in `current`.
	resultIsArray := false

	for _, c := range j.components {
		if resultIsArray {
			// replace each item in `current` by its extracted child item
			arr, err := current.asArray()
			if err != nil {
				return nil, err
			}
			for i, currentElem := range arr {
				err := c.extract(currentElem, &next)
				if err != nil {
					return nil, err
				}
				if c.resultMultiplicity() == many && next.Type() == TypeArray {
					// if we get an nil array result, turn it into an empty array instead
					if a, _ := next.asArray(); a == nil {
						next = Array{}
					}
				}
				// we assign a new value to a position of `current` here.
				// this is only valid (and does not change the input Map)
				// because all functions with `resultMultiplicity() == many`
				// are required to return a *new* slice, not a pointer
				// to an existing one!
				arr[i] = next

			}
		} else {
			// replace `current` by its extracted child item
			err := c.extract(current, &next)
			if err != nil {
				return nil, err
			}
			if c.resultMultiplicity() == many && next.Type() == TypeArray {
				// if we get an nil array result, turn it into an empty array instead
				if a, _ := next.asArray(); a == nil {
					next = Array{}
				}
			}
			current = next
		}
		resultIsArray = resultIsArray || (c.resultMultiplicity() == many)
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

// extractor describes an entity that can extract a child element
// from a Value.
type extractor interface {
	extract(v Value, next *Value) error
	extractForSet(Value, *Value, *func(Value)) error
	resultMultiplicity() multiplicity
}

// addMapAccess is called when we discover `foo` or `['bar']`
// in a JSON Path string.
func (j *jsonPeg) addMapAccess(s string) {
	j.components = append(j.components, &mapValueExtractor{s})
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

func (a *mapValueExtractor) resultMultiplicity() multiplicity {
	return one
}

// addRecursiveAccess is called when we discover `..foo` or `..['bar']`
// in a JSON Path string.
func (j *jsonPeg) addRecursiveAccess(s string) {
	j.components = append(j.components, &recursiveExtractor{s})
}

// recursiveExtractor can extract a list of all items with a certain key,
// no matter where they are located in the Map
type recursiveExtractor struct {
	key string
}

func (a *recursiveExtractor) extract(v Value, next *Value) error {
	var results []Value

	if v.Type() == TypeMap {
		// if v is a Map, then we append the entry with the correct
		// key (if one exists) to the result list and recurse for all
		// contained containers
		cont, _ := v.asMap()
		for key, value := range cont {
			if key == a.key {
				// NB. We do NOT descend further into `value` even if
				// it is itself a Map or Array!
				results = append(results, value)
			} else if value.Type() == TypeMap || value.Type() == TypeArray {
				// recurse
				var descend Value
				err := a.extract(value, &descend)
				if err != nil {
					return err
				}
				// we expect that the results we get from further
				// down are in an array shape
				subResults, err := descend.asArray()
				if err != nil {
					return err
				}
				results = append(results, subResults...)
			}
			// we ignore all entries in this Map that are not
			// container-like or do not have the key we are looking for
		}
	} else if v.Type() == TypeArray {
		// if v is a Map, then we
		cont, _ := v.asArray()
		for _, value := range cont {
			if value.Type() == TypeMap || value.Type() == TypeArray {
				// recurse
				var descend Value
				err := a.extract(value, &descend)
				if err != nil {
					return err
				}
				// we expect that the results we get from further
				// down are in an array shape
				subResults, err := descend.asArray()
				if err != nil {
					return err
				}
				results = append(results, subResults...)
			}
			// we ignore all entries in this Array that are not
			// container-like
		}
	} else {
		return fmt.Errorf("cannot descend recursively into %T", v)
	}

	*next = Array(results)
	return nil
}

func (a *recursiveExtractor) extractForSet(v Value, next *Value, setInParent *func(Value)) error {
	return fmt.Errorf("not implemented")
}

func (a *recursiveExtractor) resultMultiplicity() multiplicity {
	return many
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

func (a *arrayElementExtractor) resultMultiplicity() multiplicity {
	return one
}

// addArraySlice is called when we discover `[1:3]` or `[1:3:2]` in a
// JSON Path string.
func (j *jsonPeg) addArraySlice(s string) {
	parts := strings.Split(s, ":")
	if !(len(parts) == 2 || len(parts) == 3) {
		panic(fmt.Sprintf("'%s' did not have format 'a:b' or 'a:b:c'", s))
	}
	// due to parser configuration, s will always contain numeric strings,
	// but they may overflow int32, so we need a check here.
	start, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		panic(fmt.Sprintf("overflow index number: " + parts[0]))
	}
	end, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		panic(fmt.Sprintf("overflow index number: " + parts[1]))
	}
	step := int64(1)
	if len(parts) == 3 {
		step, err = strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			panic(fmt.Sprintf("overflow index number: " + parts[2]))
		}
	}
	if start > end {
		panic(fmt.Sprintf("start index %d must be less or equal to end index %d",
			start, end))
	}
	j.components = append(j.components, &arraySliceExtractor{int(start), int(end), int(step)})
}

// arraySliceExtractor can extract a slice from an Array using the
// given start/end indexes.
type arraySliceExtractor struct {
	start, end, step int
}

func (a *arraySliceExtractor) extract(v Value, next *Value) error {
	cont, err := AsArray(v)
	if err != nil {
		return fmt.Errorf("cannot access a %T using range %d:%d", v, a.start, a.end)
	}
	// negative indexes are forbidden at the moment
	if a.start < 0 || a.end < 0 {
		return fmt.Errorf("array indexes must be >= 0")
	}
	// end index must be greater or equal than start index
	if a.start > a.end {
		return fmt.Errorf("start index %d must be less or equal to end index %d",
			a.start, a.end)
	}
	if a.start >= len(cont) {
		*next = Array{}
		return nil
	}
	// if too large, truncate the end index to the largest possible value
	end := a.end
	if a.end > len(cont) {
		end = len(cont)
	}

	// copy the values into a new array
	retVal := make(Array, 0, end-a.start)
	for i := a.start; i < end; i += a.step {
		retVal = append(retVal, cont[i])
	}
	*next = retVal
	return nil
}

func (a *arraySliceExtractor) extractForSet(v Value, next *Value, setInParent *func(Value)) error {
	return fmt.Errorf("not implemented")
}

func (a *arraySliceExtractor) resultMultiplicity() multiplicity {
	return many
}
