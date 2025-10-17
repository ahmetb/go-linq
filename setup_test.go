package linq

import (
	"fmt"
	"iter"
	"slices"
	"testing"
)

type foo struct {
	f1 int
	f2 bool
	f3 string
}

func (f foo) Iterate() iter.Seq[any] {
	return func(yield func(any) bool) {
		// Yield the first field and check if we should continue.
		if !yield(f.f1) {
			return
		}

		// Yield the second field and check if we should continue.
		if !yield(f.f2) {
			return
		}

		// Yield the third and final field.
		yield(f.f3)
	}
}

func (f foo) CompareTo(c Comparable) int {
	a, b := f.f1, c.(foo).f1

	if a < b {
		return -1
	} else if a > b {
		return 1
	}

	return 0
}

func toSlice(q Query) (result []any) {
	q.Iterate(func(item any) bool {
		result = append(result, item)
		return true
	})

	return
}

// testQueryIteration tests the iteration of a query. First, it aborts the
// iteration by returning false. Then, it verifies that the output of the
// iteration is as expected.
//
// NOTE: This function might not behave as expected if the query does not
// support reiteration, e.g., iteration over a channel.
func testQueryIteration(q Query, expected []any) bool {
	q.Iterate(func(item any) bool { return false })
	return verifyQueryOutput(q, expected)
}

// verifyQueryOutput verifies that the output of a query is as expected.
func verifyQueryOutput(q Query, expected []any) (result bool) {
	actual := toSlice(q)
	result = slices.Equal(actual, expected)
	if !result {
		fmt.Printf("got=[%v] expected=[%v]", actual, expected)
	}
	return
}

func mustPanicWithError(t *testing.T, expectedErr string, f func()) {
	defer func() {
		r := recover()
		err := fmt.Sprintf("%s", r)
		if err != expectedErr {
			t.Fatalf("got=[%v] expected=[%v]", err, expectedErr)
		}
	}()
	f()
}
