package linq

// IntersectBy yields the first value for each key occurring in both queries.
func (q Query[T]) IntersectBy[K comparable](other Query[T], selector func(T) K) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		available := make(map[K]struct{})
		for value := range other.iterate {
			available[selector(value)] = struct{}{}
		}

		q.iterate(func(value T) bool {
			key := selector(value)
			if _, found := available[key]; !found {
				return true
			}

			delete(available, key)
			return yield(value)
		})
	}}
}
