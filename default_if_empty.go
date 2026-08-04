package linq

// DefaultIfEmpty yields defaultValue only when the source has no values.
func (q Query[T]) DefaultIfEmpty(defaultValue T) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
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
