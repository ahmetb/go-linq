package linq

// Prepend returns a lazy query that yields value before the source.
func (q Query[T]) Prepend(value T) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		if yield(value) {
			q.iterate(yield)
		}
	}}
}
