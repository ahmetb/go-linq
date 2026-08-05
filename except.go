package linq

// Except produces the set difference of two sequences. The set difference is
// the members of the first sequence that don't appear in the second sequence.
//
// Elements of basic comparable kinds (integers, floats, complex numbers,
// strings, booleans) are tracked in a strongly-typed set with no boxing. All
// other element types are tracked by their boxed (interface) values, so for
// them this method panics if T is not a comparable type at runtime; ExceptBy
// with a comparable key selector is the fast path for such types.
func (q Query[T]) Except(q2 Query[T]) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := newSeenSet[T]()
			q2.Iterate(func(item T) bool {
				set.add(item)
				return true
			})

			q.Iterate(func(item T) bool {
				if !set.has(item) {
					return yield(item)
				}
				return true
			})
		},
	}
}

// ExceptBy invokes a transform function on each element of a collection and
// produces the set difference of two sequences. The set difference is the
// members of the first sequence that don't appear in the second sequence.
//
// ExceptBy is a generic method: the comparison key type TKey is inferred from
// the selector function and must be comparable.
func (q Query[T]) ExceptBy[TKey comparable](q2 Query[T], selector func(T) TKey) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[TKey]struct{})
			q2.Iterate(func(item T) bool {
				set[selector(item)] = struct{}{}
				return true
			})

			q.Iterate(func(item T) bool {
				if _, seen := set[selector(item)]; !seen {
					return yield(item)
				}
				return true
			})
		},
	}
}
