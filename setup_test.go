package linq

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

type foo struct {
	f1 int
	f2 bool
	f3 string
}

func toSlice[T any](q Query[T]) (result []T) {
	q.Iterate(func(item T) bool {
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
func testQueryIteration[T any](q Query[T], expected []T) bool {
	runDryIteration(q)
	return assertQueryOutput(q, expected)
}

// runDryIteration performs a no-op iteration over the query
// to test whether it supports early abort and reiteration.
func runDryIteration[T any](q Query[T]) {
	q.Iterate(func(item T) bool { return false })
}

// assertQueryOutput verifies that the output of a query is as expected.
func assertQueryOutput[T any](q Query[T], expected []T) (result bool) {
	actual := toSlice(q)
	result = slices.EqualFunc(actual, expected, func(a, b T) bool {
		return reflect.DeepEqual(a, b)
	})
	if !result {
		fmt.Printf("got=[%v] expected=[%v]", actual, expected)
	}
	return
}

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	f()
}
