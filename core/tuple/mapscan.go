package tuple

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"unicode"
)

var reArrayPath = regexp.MustCompile(`^([^\[]+)?(\[[0-9]+\])?$`)

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
			i++
			if i < length {
				part += string(runes[i])
			}
		case '.':
			if part != "" {
				result = append(result, part)
				part = ""
			}
		case '[':
			if i < length-1 {
				nr := runes[i+1]
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

// get inner string, if non-sense brackets then return empty
func splitBracket(runes []rune, i int, quote rune) string {
	length := len(runes)
	result := ""
	for i < length {
		r := runes[i]
		if r == quote {
			if i < length-1 {
				if runes[i+1] == ']' {
					index := ""
					if i < length-4 && runes[i+2] == '[' {
						index = getArrayIndex(runes, i+3)
					}
					return result + index
				}
			}
		}
		result += string(r)

		i++
	}
	return ""
}

// get array index string, if not has index then return empty
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

func scanMap(m Map, p string, v *Value) (err error) {
	if p == "" {
		return errors.New("empty key is not supported")
	}
	var tempValue Value
	tempMap := m
	for _, token := range split(p) {
		matchStr := reArrayPath.FindAllStringSubmatch(token, -1)
		if len(matchStr) == 0 {
			return errors.New("invalid path phrase")
		}
		submatchStr := matchStr[0]
		if submatchStr[1] != "" {
			foundMap := tempMap[submatchStr[1]]
			if foundMap == nil {
				return errors.New(
					"not found the key in map: " + submatchStr[1])
			}
			tempValue = foundMap
			if foundMap.Type() == TypeMap {
				tempMap, _ = foundMap.Map()
			}
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
			a, err := tempValue.Array()
			if err != nil {
				return err
			}
			i := int(i64)
			if i >= len(a) {
				return errors.New("out of range access: " + token)
			}
			tempValue = a[i]
		}
	}
	*v = tempValue
	return nil
}
