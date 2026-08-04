package linq

func (q query[T]) Join[U any, K comparable, R any](
	inner query[U],
	outerKeySelector func(T) K,
	innerKeySelector func(U) K,
	resultSelector func(T, U) R,
) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			lookup := make(map[K][]U)
			for value := range inner.iterate {
				key := innerKeySelector(value)
				lookup[key] = append(lookup[key], value)
			}

			q.iterate(func(outer T) bool {
				for _, matched := range lookup[outerKeySelector(outer)] {
					if !yield(resultSelector(outer, matched)) {
						return false
					}
				}
				return true
			})
		},
	}
}
