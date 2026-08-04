package linq

func (q Query[T]) Prepend(value T) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		if yield(value) {
			q.iterate(yield)
		}
	}}
}
