package linq

func (q query[T]) GroupJoin[U any, K comparable, R any](
	inner query[U],
	outerKeySelector func(T) K,
	innerKeySelector func(U) K,
	resultSelector func(T, []U) R,
) query[R] {
	return query[R]{
		iterate: func(yield func(R) bool) {
			lookup := make(map[K][]U)
			for value := range inner.iterate {
				key := innerKeySelector(value)
				lookup[key] = append(lookup[key], value)
			}
			empty := []U{}

			q.iterate(func(outer T) bool {
				matched, ok := lookup[outerKeySelector(outer)]
				if !ok {
					matched = empty
				}
				return yield(resultSelector(outer, matched))
			})
		},
	}
}
