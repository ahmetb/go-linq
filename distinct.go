package linq

// Distinct method returns distinct elements from a collection. The result is an
// unordered collection that contains no duplicate values.
//
// Elements are tracked in a set keyed by their boxed (interface) values, so
// this method panics if T is not a comparable type at runtime. For element
// types whose boxing allocates (strings, structs, large numbers), DistinctBy
// with an identity selector avoids the boxing and performs better.
func (q Query[T]) Distinct() Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[any]struct{})

			q.Iterate(func(item T) bool {
				if _, seen := set[item]; !seen {
					set[item] = struct{}{}
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
			var previous T
			isFirst := true

			oq.Iterate(func(item T) bool {
				if isFirst || any(item) != any(previous) {
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
