package execution

import (
	"pfi/sensorbee/sensorbee/data"
)

// sortArray holds an array of values of which we want to find the order,
// plus a flag whether we want ascending order or not.
type sortArray struct {
	values    data.Array
	ascending bool
}

// indexSlice is a data structure that allows to put the `indexes` integer
// array in the order of items in the `ordering` arrays.
//
// If `indexes` holds the values {0, 1, 2} and `ordering` is
//   []sortArray{
//     {[]data.Value{data.Int(1), data.Int(1), data.Int(1)}, true},
//     {[]data.Value{data.Int(5), data.Int(5), data.Int(6)}, false},
//     {[]data.Value{data.Int(8), data.Int(9), data.Int(7)}, true},
//   }
// then after sort.Sort() on that indexSlice, `indexes` will have
// the value `{2, 0, 1}` because `ordering[2]` is the smallest item
// according to the sort criteria, then `ordering[0]`, then `ordering[1]`.
type indexSlice struct {
	indexes  []int
	ordering []sortArray
}

func (s *indexSlice) Len() int {
	return len(s.indexes)
}

func (s *indexSlice) Swap(i, j int) {
	s.indexes[i], s.indexes[j] = s.indexes[j], s.indexes[i]
}

func (s *indexSlice) Less(i, j int) bool {
	for _, order := range s.ordering {
		iVal := order.values[s.indexes[i]]
		jVal := order.values[s.indexes[j]]
		if data.Equal(iVal, jVal) {
			continue
		} else if order.ascending {
			return data.Less(iVal, jVal)
		}
		return data.Less(jVal, iVal)
	}
	return false
}
