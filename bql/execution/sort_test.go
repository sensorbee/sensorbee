package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"sort"
	"testing"
)

func TestSorting(t *testing.T) {
	testCases := []struct {
		ordering []sortArray
		expected []int // the order in which to traverse the array to go from small to big
	}{
		/// one sorting dimension
		// already sorted
		{[]sortArray{
			{[]data.Value{data.Int(1), data.Int(2), data.Int(3)}, true},
		}, []int{0, 1, 2}},
		// some other order
		{[]sortArray{
			{[]data.Value{data.Int(2), data.Int(3), data.Int(1)}, true},
		}, []int{2, 0, 1}},
		// sort descending
		{[]sortArray{
			{[]data.Value{data.Int(1), data.Int(2), data.Int(3)}, false},
		}, []int{2, 1, 0}},
		// sort some other order descending
		{[]sortArray{
			{[]data.Value{data.Int(2), data.Int(3), data.Int(1)}, false},
		}, []int{1, 0, 2}},
		/// two sorting dimensions
		// already sorted
		{[]sortArray{
			{[]data.Value{data.Int(1), data.Int(2), data.Int(3)}, true},
			{[]data.Value{data.Int(4), data.Int(5), data.Int(6)}, true},
		}, []int{0, 1, 2}},
		// second is ignored if first is correct
		{[]sortArray{
			{[]data.Value{data.Int(1), data.Int(2), data.Int(3)}, true},
			{[]data.Value{data.Int(4), data.Int(5), data.Int(6)}, false},
		}, []int{0, 1, 2}},
		// second is used if first is the same everywhere
		{[]sortArray{
			{[]data.Value{data.Int(1), data.Int(1), data.Int(1)}, true},
			{[]data.Value{data.Int(4), data.Int(5), data.Int(6)}, false},
		}, []int{2, 1, 0}},
		// third is used if first and second are the same
		{[]sortArray{
			{[]data.Value{data.Int(1), data.Int(1), data.Int(1)}, true},
			{[]data.Value{data.Int(5), data.Int(5), data.Int(6)}, false},
			{[]data.Value{data.Int(8), data.Int(9), data.Int(7)}, true},
		}, []int{2, 0, 1}},
	}

	for _, tc := range testCases {
		ordering := tc.ordering
		expected := tc.expected
		Convey(fmt.Sprintf("Given sort input %v", ordering), t, func() {
			indexes := make([]int, len(tc.ordering[0].values))
			for i := range tc.ordering[0].values {
				indexes[i] = i
			}
			is := &indexSlice{indexes, ordering}

			Convey("Then the sort output should match", func() {
				sort.Sort(is)
				So(is.indexes, ShouldResemble, expected)
			})
		})
	}
}
