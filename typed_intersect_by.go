package linq

func (q query[T]) IntersectBy[K comparable](other query[T], selector func(T) K) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
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
