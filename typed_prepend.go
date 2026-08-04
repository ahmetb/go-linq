package linq

func (q query[T]) Prepend(value T) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		if yield(value) {
			q.iterate(yield)
		}
	}}
}
