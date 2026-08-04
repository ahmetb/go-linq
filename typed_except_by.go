package linq

func (q Query[T]) ExceptBy[K comparable](other Query[T], selector func(T) K) Query[T] {
	return Query[T]{
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
