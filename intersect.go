package linq

// Intersect produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
//
// Elements are tracked in a set keyed by their boxed (interface) values, so
// this method panics if T is not a comparable type at runtime. IntersectBy
// with an identity selector avoids the boxing and performs better.
func (q Query[T]) Intersect(q2 Query[T]) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[any]struct{})
			for item := range q2.Iterate {
				set[item] = struct{}{}
			}

			for item := range q.Iterate {
				if _, exists := set[item]; exists {
					delete(set, item)
					if !yield(item) {
						return
					}
				}
			}
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
			for item := range q2.Iterate {
				set[selector(item)] = struct{}{}
			}

			for item := range q.Iterate {
				key := selector(item)
				if _, exists := set[key]; exists {
					delete(set, key)
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}
