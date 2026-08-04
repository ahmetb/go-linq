package linq

// Concat returns a lazy query that yields q followed by other.
func (q Query[T]) Concat(other Query[T]) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		continuing := true
		q.iterate(func(value T) bool {
			continuing = yield(value)
			return continuing
		})
		if continuing {
			other.iterate(yield)
		}
	}}
}
