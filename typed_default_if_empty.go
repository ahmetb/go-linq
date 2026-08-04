package linq

func (q query[T]) DefaultIfEmpty(defaultValue T) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		seen := false
		continuing := true
		q.iterate(func(value T) bool {
			seen = true
			continuing = yield(value)
			return continuing
		})
		if !seen && continuing {
			yield(defaultValue)
		}
	}}
}
