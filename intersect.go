package linq

// Intersect produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
//
// Elements of basic comparable kinds (integers, floats, complex numbers,
// strings, booleans) are tracked in a strongly-typed set with no boxing. All
// other element types are tracked by their boxed (interface) values, so for
// them this method panics if T is not a comparable type at runtime;
// IntersectBy with a comparable key selector is the fast path for such types.
func (q Query[T]) Intersect(q2 Query[T]) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := newSeenSet[T]()
			q2.Iterate(func(item T) bool {
				set.add(item)
				return true
			})

			q.Iterate(func(item T) bool {
				if set.has(item) {
					set.del(item)
					return yield(item)
				}
				return true
			})
		},
	}
}

// IntersectBy produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
//
// IntersectBy invokes a transform function on each element of both collections.
// It is a generic method: the comparison key type TKey is inferred from the
// selector function and must be comparable.
func (q Query[T]) IntersectBy[TKey comparable](q2 Query[T], selector func(T) TKey) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[TKey]struct{})
			q2.Iterate(func(item T) bool {
				set[selector(item)] = struct{}{}
				return true
			})

			q.Iterate(func(item T) bool {
				key := selector(item)
				if _, exists := set[key]; exists {
					delete(set, key)
					return yield(item)
				}
				return true
			})
		},
	}
}
