package data

// SummarizeConfig (will) has cnofiguration parameters of Summarize function.
type SummarizeConfig struct {
}

// Summarize summarizes a Value as a JSON-like string. It may truncate or
// filter out the content of strings and blobs. The result may not be able
// to parsed as a JSON.
//
// It'll support the max depth of a map, the max number of fields of a map
// to be rendered, or the max number of elements in an array, and so on.
func Summarize(val Value) string {
	v := val.clone()
	return summarize(v).String()
}

func summarize(v Value) Value {
	switch v.Type() {
	case TypeBlob:
		return String("(blob)")

	case TypeMap:
		m, _ := v.asMap()
		for k, val := range m {
			m[k] = summarize(val)
		}

	case TypeArray:
		a, _ := v.asArray()
		for i, val := range a {
			a[i] = summarize(val)
		}
	}
	return v
}
