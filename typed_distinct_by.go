package linq

func (q query[T]) DistinctBy[K comparable](selector func(T) K) query[T] {
	return query[T]{
		iterate: func(yield func(T) bool) {
			seen := make(map[K]struct{})
			q.iterate(func(value T) bool {
				key := selector(value)
				if _, ok := seen[key]; ok {
					return true
				}
				seen[key] = struct{}{}
				return yield(value)
			})
		},
	}
}
