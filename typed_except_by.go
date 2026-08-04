package linq

func (q query[T]) ExceptBy[K comparable](other query[T], selector func(T) K) query[T] {
	return query[T]{
		iterate: func(yield func(T) bool) {
			excluded := make(map[K]struct{})
			for value := range other.iterate {
				excluded[selector(value)] = struct{}{}
			}

			q.iterate(func(value T) bool {
				key := selector(value)
				if _, found := excluded[key]; found {
					return true
				}
				excluded[key] = struct{}{}
				return yield(value)
			})
		},
	}
}
