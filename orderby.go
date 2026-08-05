package linq

import (
	"cmp"
	"iter"
	"slices"
	"sort"
)

// OrderedQuery is the type returned from OrderBy, OrderByDescending ThenBy and
// ThenByDescending functions.
type OrderedQuery[T any] struct {
	Query[T]
	original Query[T]
	compares []func(a, b T) int
}

func ascending[T any, TKey cmp.Ordered](selector func(T) TKey) func(a, b T) int {
	return func(a, b T) int {
		return cmp.Compare(selector(a), selector(b))
	}
}

func descending[T any, TKey cmp.Ordered](selector func(T) TKey) func(a, b T) int {
	return func(a, b T) int {
		return cmp.Compare(selector(b), selector(a))
	}
}

// sortedIterate returns a sequence that collects the query into a slice and
// sorts it with the given comparison functions, applied in order until one of
// them reports a difference.
func (q Query[T]) sortedIterate(compares []func(a, b T) int) iter.Seq[T] {
	compare := compares[0]
	if len(compares) > 1 {
		compare = func(a, b T) int {
			for _, compare := range compares {
				if c := compare(a, b); c != 0 {
					return c
				}
			}
			return 0
		}
	}

	return func(yield func(T) bool) {
		items := q.collect()

		slices.SortFunc(items, compare)

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}

// OrderBy sorts the elements of a collection in ascending order. Elements are
// sorted according to a key.
//
// OrderBy is a generic method: the key type TKey is inferred from the selector
// function and must be an ordered type (cmp.Ordered). To sort by a custom
// comparison instead, use the Sort method.
func (q Query[T]) OrderBy[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	compares := []func(a, b T) int{ascending(selector)}
	return OrderedQuery[T]{
		compares: compares,
		original: q,
		Query: Query[T]{
			Iterate: q.sortedIterate(compares),
			size:    q.size,
		},
	}
}

// OrderByDescending sorts the elements of a collection in descending order.
// Elements are sorted according to a key.
//
// OrderByDescending is a generic method: the key type TKey is inferred from the
// selector function and must be an ordered type (cmp.Ordered).
func (q Query[T]) OrderByDescending[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	compares := []func(a, b T) int{descending(selector)}
	return OrderedQuery[T]{
		compares: compares,
		original: q,
		Query: Query[T]{
			Iterate: q.sortedIterate(compares),
			size:    q.size,
		},
	}
}

// ThenBy performs a subsequent ordering of the elements in a collection in
// ascending order. This method enables you to specify multiple sort criteria by
// applying any number of ThenBy or ThenByDescending methods.
func (oq OrderedQuery[T]) ThenBy[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	compares := append(slices.Clip(oq.compares), ascending(selector))
	return OrderedQuery[T]{
		compares: compares,
		original: oq.original,
		Query: Query[T]{
			Iterate: oq.original.sortedIterate(compares),
			size:    oq.original.size,
		},
	}
}

// ThenByDescending performs a subsequent ordering of the elements in a
// collection in descending order. This method enables you to specify multiple
// sort criteria by applying any number of ThenBy or ThenByDescending methods.
func (oq OrderedQuery[T]) ThenByDescending[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	compares := append(slices.Clip(oq.compares), descending(selector))
	return OrderedQuery[T]{
		compares: compares,
		original: oq.original,
		Query: Query[T]{
			Iterate: oq.original.sortedIterate(compares),
			size:    oq.original.size,
		},
	}
}

// sorter adapts a less function to sort.Interface so that sorting calls the
// user's comparator exactly once per comparison.
type sorter[T any] struct {
	items []T
	less  func(i, j T) bool
}

func (s sorter[T]) Len() int           { return len(s.items) }
func (s sorter[T]) Swap(i, j int)      { s.items[i], s.items[j] = s.items[j], s.items[i] }
func (s sorter[T]) Less(i, j int) bool { return s.less(s.items[i], s.items[j]) }

// Sort returns a new query by sorting elements with provided less function in
// ascending order. The comparer function should return true if the parameter i
// is less than j.
//
// Unlike OrderBy, Sort does not require the sort key to be an ordered type,
// so it can be used with arbitrary comparison logic. The less function is
// invoked once per element comparison.
func (q Query[T]) Sort(less func(i, j T) bool) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			items := q.collect()

			sort.Sort(sorter[T]{items: items, less: less})

			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		},
		size: q.size,
	}
}
