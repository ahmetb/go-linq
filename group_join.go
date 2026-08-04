package linq

// GroupJoin correlates each outer value with all inner values sharing its key.
func (q Query[T]) GroupJoin[U any, K comparable, R any](
	inner Query[U],
	outerKeySelector func(T) K,
	innerKeySelector func(U) K,
	resultSelector func(T, []U) R,
) Query[R] {
	return Query[R]{
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
