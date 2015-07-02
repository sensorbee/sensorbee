package data

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"unicode"
)

var reArrayPath = regexp.MustCompile(`^([^\[]+)?(\[[0-9]+\])?$`)

// split splits a string describing a JSON path into its components,
// i.e., strings representing one level of descend in a map. Array
// indexes in brackets are returned together with their "parent string".
//
// Examples:
//  split(`["store"]["book"][0]["title"]`)
//  // []string{"store", "book[0]", "title"}
//  split(`store.book[0].title`)
//  // []string{"store", "book[0]", "title"}
func split(s string) []string {
	i := 0
	result := []string{}
	part := ""
	runes := []rune(s)
	length := len(runes)
	for i < length {
		r := runes[i]
		switch r {
		case '\\':
			// if we see a backslash, the next character will be
			// appended to the string verbatim (i.e. we can escape
			// dots, brackets etc outside of brackets). if it is
			// the last character of the input string it will be ignored.
			i++
			if i < length {
				part += string(runes[i])
			}
		case '.':
			// if we see a dot, we will finish the current component
			// and append it to the result list. if the current
			// component is empty, the dot will be ignored.
			if part != "" {
				result = append(result, part)
				part = ""
			}
		case '[':
			// if we see an opening bracket, we can either have
			// hoge[123] or hoge["key"] or some invalid situation.
			if i < length-1 {
				nr := runes[i+1]
				// if we have the hoge["key"] situation, get the
				// part until the closing bracket.
				if nr == '"' || nr == '\'' {
					inbracket := splitBracket(runes, i+2, nr)
					if inbracket != "" {
						if part != "" {
							result = append(result, part)
						}
						result = append(result, inbracket)
						part = ""
						i += 1 + len(inbracket) + 2 // " + inner bracket + "]
						break
					}
				}
			}
			// NB. if we have a string after an opening bracket that
			// does not begin with a quote character (this can be numeric
			// or not), it will be copied verbatim
			part += string(r)
		default:
			part += string(r)
		}
		i++
	}
	if part != "" {
		result = append(result, part)
	}
	return result
}

// splitBracket returns a string in `runes` that begins at position `i`
// and is followed by the given quote character and a closing bracket.
// It can be used to extract `key` from `hoge["key"]`. If there is an
// integer index in brackets following, this will be returned as well.
// If there is no string matching the conditions, returns an empty string.
//
// Example:
//  splitBracket([]rune(`a["hoge"].b`), 3, '"')
//  // `hoge`
//  splitBracket([]rune(`a["hoge"][123]`), 3, '"')
//  // `hoge[123]`
func splitBracket(runes []rune, i int, quote rune) string {
	length := len(runes)
	result := ""
	for i < length {
		r := runes[i]
		// if the current character is the required quote character ...
		if r == quote {
			if i < length-1 {
				// ... and the next character is the closing bracket:
				if runes[i+1] == ']' {
					index := ""
					// if there is a following opening bracket, which
					// may or may not be followed by an array index,
					// try to get that index and append it to what
					// we found so far (it will be an empty string if
					// there is no array index but a string found)
					if i < length-4 && runes[i+2] == '[' {
						index = getArrayIndex(runes, i+3)
					}
					// return what we found until now
					return result + index
				}
			}
		}
		// otherwise just append the found character to the intermediate
		// result
		result += string(r)

		i++
	}
	// if we never fulfill the above condition, return an empty string
	return ""
}

// getArrayIndex returns a string in `runes` that begins at position `i`,
// consists only of digits and is followed by a closing bracket.
// The string that was found is wrapped in brackets before returning.
// It can be used to extract `[123]` from `hoge[123]`. If there is no
// string matching the conditions, returns an empty string.
//
// Example:
//  getArrayIndex([]rune(`hoge[123]`), 5)
//  // `[123]`
func getArrayIndex(runes []rune, i int) string {
	length := len(runes)
	result := ""
	for i < length {
		r := runes[i]
		if r == ']' {
			break
		} else if unicode.IsNumber(r) {
			result += string(r)
		} else {
			return ""
		}
		i++
	}
	if result != "" {
		return "[" + result + "]"
	}
	return ""
}

// scanMap does basically what is described in the Map.Get documentation.
// The value found at p is written to v.
func scanMap(m Map, p string, v *Value) (err error) {
	if p == "" {
		return errors.New("empty key is not supported")
	}
	// tempValue will point to the item of the Map that we are
	// currently investigating
	var tempValue Value = m
	// loop over the components of the path, like "key" or "hoge[123]"
	for _, token := range split(p) {
		// check that we do indeed have a valid component form
		matchStr := reArrayPath.FindAllStringSubmatch(token, -1)
		if len(matchStr) == 0 {
			return errors.New("invalid path component: " + token)
		}
		// get the "before brackets" part of the component
		submatchStr := matchStr[0]
		if submatchStr[1] != "" {
			// try to access the current tempValue as a map and
			// pull out the value therein
			tempMap, err := tempValue.asMap()
			if err != nil {
				return fmt.Errorf("cannot access a %T using key \"%s\"",
					tempValue, token)
			}
			foundValue := tempMap[submatchStr[1]]
			if foundValue == nil {
				return errors.New(
					"not found the key in map: " + submatchStr[1])
			}
			tempValue = foundValue
		}
		// get array index number
		if submatchStr[2] != "" {
			i64, err := strconv.ParseInt(
				submatchStr[2][1:len(submatchStr[2])-1], 10, 64)
			if err != nil {
				return errors.New("invalid array index number: " + token)
			}
			if i64 > math.MaxInt32 {
				return errors.New("overflow index number: " + token)
			}
			// try to access the current tempValue as an array
			// and access the value therein
			tempArr, err := tempValue.asArray()
			if err != nil {
				return fmt.Errorf("cannot access a %T using index %d",
					tempValue, i64)
			}
			i := int(i64)
			if i >= len(tempArr) {
				return errors.New("out of range access: " + token)
			}
			tempValue = tempArr[i]
		}
	}
	*v = tempValue
	return nil
}

// setInMap does basically what is described in the Map.Set documentation.
// The value at v is written to m at the given path.
func setInMap(m Map, p string, v Value) (err error) {
	if p == "" {
		return errors.New("empty key is not supported")
	}
	// tempValue will point to the item of the Map that we are
	// currently investigating
	var tempValue Value = m
	var setValueAtCurrentPos func(Value)
	// loop over the components of the path, like "key" or "hoge[123]"
	components := split(p)
	for cIdx, token := range components {
		last := cIdx == len(components)-1
		// check that we do indeed have a valid component form
		matchStr := reArrayPath.FindAllStringSubmatch(token, -1)
		if len(matchStr) == 0 {
			return errors.New("invalid path component: " + token)
		}
		// get the "before brackets" part of the component
		submatchStr := matchStr[0]
		hasKey := submatchStr[1] != ""
		hasIndex := submatchStr[2] != ""
		if hasKey {
			// normally we would fail when trying to access a null element
			// with a key (like `null["hoge"]`), but in this special situation
			// we will create a map there so that a "hoge" key can be added
			if tempValue.Type() == TypeNull {
				tempValue = Map{}
				setValueAtCurrentPos(tempValue)
			}
			// try to access the current tempValue as a map and
			// pull out the value therein
			tempMap, err := tempValue.asMap()
			if err != nil {
				return fmt.Errorf("cannot access a %T using key \"%s\"",
					tempValue, token)
			}
			// we need a helper function to write the value at the
			// current memory position. we do not need it here because
			// we have access to `tempMap`, but we need it below when
			// dealing with lists, where we do not have access to
			// `tempMap` any more and otherwise could only do changes
			// in-place.
			k := submatchStr[1]
			setValueAtCurrentPos = func(v Value) {
				tempMap[k] = v
			}

			if last && !hasIndex {
				// if this is already the last component, we need
				// to set the value in the map
				setValueAtCurrentPos(v)
			} else {
				// otherwise, we will have to go one level deeper
				foundValue := tempMap[submatchStr[1]]
				if foundValue == nil {
					// create a map for further nesting (this is the
					// `mkdir -p` behavior)
					if !hasIndex {
						setValueAtCurrentPos(Map{})
					} else {
						setValueAtCurrentPos(Array{})
					}
					foundValue = tempMap[submatchStr[1]]
				}
				tempValue = foundValue
			}
		}
		// get array index number
		if hasIndex {
			i64, err := strconv.ParseInt(
				submatchStr[2][1:len(submatchStr[2])-1], 10, 64)
			if err != nil {
				return errors.New("invalid array index number: " + token)
			}
			if i64 > math.MaxInt32 {
				return errors.New("overflow index number: " + token)
			}
			// normally we would fail when trying to access a null element
			// with an index (like `null[5]`), but in this special situation
			// we will create an array there so that elements can be appended
			if tempValue.Type() == TypeNull {
				tempValue = Array{}
				setValueAtCurrentPos(tempValue)
			}
			// try to access the current tempValue as an array
			// and access the value therein
			tempArr, err := tempValue.asArray()
			if err != nil {
				return fmt.Errorf("cannot access a %T using index %d",
					tempValue, i64)
			}
			i := int(i64)
			// if the array is not long enough ...
			if i >= len(tempArr) {
				// ... we reallocate and pad with zeros
				extArr := make(Array, i+1)
				for k := range tempArr {
					extArr[k] = tempArr[k]
				}
				for k := len(tempArr); k <= i; k++ {
					extArr[k] = Null{}
				}
				setValueAtCurrentPos(extArr)
				// update the function to set the item at the current list index
				setValueAtCurrentPos = func(v Value) {
					tempArr[i] = v
				}
				tempArr = extArr
			}
			if last {
				// if this is already the last component, we need
				// to set the value in the array
				tempArr[i] = v
			} else {
				tempValue = tempArr[i]
			}
		}
	}
	return nil
}
