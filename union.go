package linq

// Union produces the set union of two collections.
//
// This method excludes duplicates from the return set. This is different
// behavior to the Concat method, which returns all the elements in the input
// collection, including duplicates.
//
// Elements of basic comparable kinds (integers, floats, complex numbers,
// strings, booleans) are tracked in a strongly-typed set with no boxing. All
// other element types are tracked by their boxed (interface) values, so for
// them this method panics if T is not a comparable type at runtime; UnionBy
// with a comparable key selector is the fast path for such types.
func (q Query[T]) Union(q2 Query[T]) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := newSeenSet[T]()
			stopped := false

			q.Iterate(func(item T) bool {
				if set.add(item) {
					if !yield(item) {
						stopped = true
						return false
					}
				}
				return true
			})

			if stopped {
				return
			}

			q2.Iterate(func(item T) bool {
				if set.add(item) {
					if !yield(item) {
						return false
					}
				}
				return true
			})
		},
	}
}

// UnionBy produces the set union of two collections according to a specified
// key selector function.
//
// UnionBy is a generic method: the comparison key type TKey is inferred from the
// selector function and must be comparable. Elements are tracked in a
// strongly-typed set, so no boxing occurs.
func (q Query[T]) UnionBy[TKey comparable](q2 Query[T], selector func(T) TKey) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[TKey]struct{})
			stopped := false

			q.Iterate(func(item T) bool {
				key := selector(item)
				if _, seen := set[key]; !seen {
					set[key] = struct{}{}
					if !yield(item) {
						stopped = true
						return false
					}
				}
				return true
			})

			if stopped {
				return
			}

			q2.Iterate(func(item T) bool {
				key := selector(item)
				if _, seen := set[key]; !seen {
					set[key] = struct{}{}
					if !yield(item) {
						return false
					}
				}
				return true
			})
		},
	}
}
