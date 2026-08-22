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
//
// The sort is stable so elements that compare equal retain their input order,
// matching .NET LINQ.
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

		slices.SortStableFunc(items, compare)

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}

// newOrderedQuery builds what every ordering operator returns: the sorted
// sequence, plus the unsorted original and the comparison chain that ThenBy
// and ThenByDescending extend.
func newOrderedQuery[T any](original Query[T], compares []func(a, b T) int) OrderedQuery[T] {
	return OrderedQuery[T]{
		compares: compares,
		original: original,
		Query: Query[T]{
			Iterate: original.sortedIterate(compares),
			size:    original.size,
		},
	}
}

// Order sorts the elements of a collection in ascending order.
func Order[T cmp.Ordered](q Query[T]) OrderedQuery[T] {
	return q.OrderBy(func(item T) T { return item })
}

// OrderBy sorts the elements of a collection in ascending order. Elements are
// sorted according to a key.
//
// The sort is stable: elements with equal keys retain their original relative
// order, matching .NET LINQ. OrderByDescending, ThenBy and ThenByDescending
// share this guarantee.
//
// OrderBy is a generic method: the key type TKey is inferred from the selector
// function and must be an ordered type (cmp.Ordered). To sort by a custom
// comparison instead, use OrderWith, or Sort when stability and chaining are
// not needed.
func (q Query[T]) OrderBy[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	return newOrderedQuery(q, []func(a, b T) int{ascending(selector)})
}

// OrderByDescending sorts the elements of a collection in descending order.
// Elements are sorted according to a key.
//
// OrderByDescending is a generic method: the key type TKey is inferred from the
// selector function and must be an ordered type (cmp.Ordered).
func (q Query[T]) OrderByDescending[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	return newOrderedQuery(q, []func(a, b T) int{descending(selector)})
}

// OrderDescending sorts the elements of a collection in descending order.
func OrderDescending[T cmp.Ordered](q Query[T]) OrderedQuery[T] {
	return q.OrderByDescending(func(item T) T { return item })
}

// OrderWith sorts the elements of a collection with the provided comparison
// function. compare follows the same convention as cmp.Compare: a negative
// result means the first value is less than the second.
//
// Like OrderBy the sort is stable, and the result is an OrderedQuery that
// ThenBy and ThenByDescending can refine further. Unlike OrderBy it places no
// constraint on the element type. Sort accepts arbitrary comparison logic too,
// but is unstable and does not compose with ThenBy.
func (q Query[T]) OrderWith(compare func(T, T) int) OrderedQuery[T] {
	return newOrderedQuery(q, []func(a, b T) int{compare})
}

// OrderDescendingWith sorts the elements of a collection in the reverse of the
// order given by compare, which follows the same convention as cmp.Compare.
func (q Query[T]) OrderDescendingWith(compare func(T, T) int) OrderedQuery[T] {
	return q.OrderWith(func(a, b T) int { return compare(b, a) })
}

// ThenBy performs a subsequent ordering of the elements in a collection in
// ascending order. This method enables you to specify multiple sort criteria by
// applying any number of ThenBy or ThenByDescending methods.
func (oq OrderedQuery[T]) ThenBy[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	return newOrderedQuery(oq.original, append(slices.Clip(oq.compares), ascending(selector)))
}

// ThenByDescending performs a subsequent ordering of the elements in a
// collection in descending order. This method enables you to specify multiple
// sort criteria by applying any number of ThenBy or ThenByDescending methods.
func (oq OrderedQuery[T]) ThenByDescending[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	return newOrderedQuery(oq.original, append(slices.Clip(oq.compares), descending(selector)))
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
//
// Also unlike OrderBy, Sort is not stable: elements that compare equal may be
// reordered.
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
