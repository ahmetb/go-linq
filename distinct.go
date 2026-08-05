package linq

// Distinct method returns distinct elements from a collection. The result is an
// unordered collection that contains no duplicate values.
//
// Elements of basic comparable kinds (integers, floats, complex numbers,
// strings, booleans) are tracked in a strongly-typed set with no boxing. All
// other element types are tracked by their boxed (interface) values, so for
// them this method panics if T is not a comparable type at runtime; DistinctBy
// with a comparable key selector is the fast path for such types.
func (q Query[T]) Distinct() Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := newSeenSet[T]()

			q.Iterate(func(item T) bool {
				if set.add(item) {
					return yield(item)
				}

				return true
			})
		},
	}
}

// Distinct method returns distinct elements from a collection. The result is an
// ordered collection that contains no duplicate values.
//
// NOTE: Distinct method on OrderedQuery type has better performance than
// Distinct method on Query type.
func (oq OrderedQuery[T]) Distinct() OrderedQuery[T] {
	distinct := Query[T]{
		Iterate: func(yield func(T) bool) {
			equal := equalFor[T]()
			var previous T
			isFirst := true

			oq.Iterate(func(item T) bool {
				if isFirst || !equal(item, previous) {
					previous = item
					isFirst = false
					return yield(item)
				}

				return true
			})
		},
	}

	return OrderedQuery[T]{
		compares: oq.compares,
		original: distinct,
		Query:    distinct,
	}
}

// DistinctBy method returns distinct elements from a collection. This method
// executes selector function for each element to determine a value to compare.
// The result is an unordered collection that contains no duplicate values.
//
// DistinctBy is a generic method: the comparison key type TKey is inferred from
// the selector function and must be comparable. Elements are tracked in a
// strongly-typed set, so no boxing occurs.
func (q Query[T]) DistinctBy[TKey comparable](selector func(T) TKey) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[TKey]struct{})

			q.Iterate(func(item T) bool {
				key := selector(item)

				if _, seen := set[key]; !seen {
					set[key] = struct{}{}
					return yield(item)
				}

				return true
			})
		},
	}
}
