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
	compares []ordering[T]
}

// An ordering either compares elements directly or prepares a typed key cache.
// Caches belong to an enumeration, never to the reusable OrderedQuery.
type ordering[T any] struct {
	compare func(T, T) int
	// The final level breaks ties by source index; earlier levels leave ties
	// for the next level. This avoids a wrapper call for single-key sorts.
	prepare func([]T, bool) func(int, int) int
}

// keyedOrdering prepares one key per item and compares cached keys in the
// requested direction. At the final ordering level, source indices break ties
// so equal keys retain their original order.
func keyedOrdering[T any, TKey cmp.Ordered](selector func(T) TKey, descending bool) ordering[T] {
	return ordering[T]{prepare: func(items []T, last bool) func(int, int) int {
		keys := make([]TKey, len(items))
		for i, item := range items {
			keys[i] = selector(item)
		}
		if descending {
			if last {
				return func(i, j int) int {
					if c := cmp.Compare(keys[j], keys[i]); c != 0 {
						return c
					}
					return cmp.Compare(i, j)
				}
			}
			return func(i, j int) int { return cmp.Compare(keys[j], keys[i]) }
		}
		if last {
			return func(i, j int) int {
				if c := cmp.Compare(keys[i], keys[j]); c != 0 {
					return c
				}
				return cmp.Compare(i, j)
			}
		}
		return func(i, j int) int { return cmp.Compare(keys[i], keys[j]) }
	}}
}

// comparer returns an index comparator for items. Keyed orderings prepare their
// cache once for this iteration; direct orderings compare the items themselves.
func (o ordering[T]) comparer(items []T, last bool) func(int, int) int {
	if o.prepare != nil {
		return o.prepare(items, last)
	}
	return func(i, j int) int { return o.compare(items[i], items[j]) }
}

// sortedIterate sorts indices into the collected source so cached keys stay
// attached to their elements. Original indices break ties, preserving stability
// without moving elements or repeatedly evaluating selectors during sorting.
func (q Query[T]) sortedIterate(orders []ordering[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		items := q.collect()
		if len(items) == 0 {
			return
		}

		// A lone element comparator needs neither key caches nor an index map.
		if len(orders) == 1 && orders[0].compare != nil {
			slices.SortStableFunc(items, orders[0].compare)
			for _, item := range items {
				if !yield(item) {
					return
				}
			}
			return
		}

		compare := orders[0].comparer(items, len(orders) == 1)
		if len(orders) > 1 {
			compares := make([]func(int, int) int, len(orders))
			compares[0] = compare
			for i := 1; i < len(orders); i++ {
				compares[i] = orders[i].comparer(items, i == len(orders)-1)
			}
			compare = func(i, j int) int {
				for _, c := range compares {
					if result := c(i, j); result != 0 {
						return result
					}
				}
				return 0
			}
		}
		indices := make([]int, len(items))
		for i := range indices {
			indices[i] = i
		}
		slices.SortFunc(indices, compare)
		for _, i := range indices {
			if !yield(items[i]) {
				return
			}
		}
	}
}

// newOrderedQuery builds what every ordering operator returns: the sorted
// sequence, plus the unsorted original and the comparison chain that ThenBy
// and ThenByDescending extend.
func newOrderedQuery[T any](original Query[T], orders []ordering[T]) OrderedQuery[T] {
	return OrderedQuery[T]{
		compares: orders,
		original: original,
		Query: Query[T]{
			Iterate: original.sortedIterate(orders),
			size:    original.size,
		},
	}
}

// Order sorts the elements of a collection in ascending order.
func Order[T cmp.Ordered](q Query[T]) OrderedQuery[T] {
	return q.OrderWith(cmp.Compare[T])
}

// OrderBy sorts the elements of a collection in ascending order. Elements are
// sorted according to a key.
//
// The sort is stable: elements with equal keys retain their original relative
// order, matching .NET LINQ. OrderByDescending, ThenBy and ThenByDescending
// share this guarantee.
// Each key selector is evaluated once per source element per iteration,
// before sorting. Cached keys are not shared between iterations.
//
// OrderBy is a generic method: the key type TKey is inferred from the selector
// function and must be an ordered type (cmp.Ordered). To sort by a custom
// comparison instead, use OrderWith, or Sort when stability and chaining are
// not needed.
func (q Query[T]) OrderBy[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	return newOrderedQuery(q, []ordering[T]{keyedOrdering(selector, false)})
}

// OrderByDescending sorts the elements of a collection in descending order.
// Elements are sorted according to a key.
//
// OrderByDescending is a generic method: the key type TKey is inferred from the
// selector function and must be an ordered type (cmp.Ordered).
func (q Query[T]) OrderByDescending[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	return newOrderedQuery(q, []ordering[T]{keyedOrdering(selector, true)})
}

// OrderDescending sorts the elements of a collection in descending order.
func OrderDescending[T cmp.Ordered](q Query[T]) OrderedQuery[T] {
	return q.OrderDescendingWith(cmp.Compare[T])
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
	return newOrderedQuery(q, []ordering[T]{{compare: compare}})
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
	return newOrderedQuery(oq.original, append(slices.Clip(oq.compares), keyedOrdering(selector, false)))
}

// ThenByDescending performs a subsequent ordering of the elements in a
// collection in descending order. This method enables you to specify multiple
// sort criteria by applying any number of ThenBy or ThenByDescending methods.
func (oq OrderedQuery[T]) ThenByDescending[TKey cmp.Ordered](selector func(T) TKey) OrderedQuery[T] {
	return newOrderedQuery(oq.original, append(slices.Clip(oq.compares), keyedOrdering(selector, true)))
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
