package linq

func (q Query[T]) Iterate(yield func(T) bool) {
	q.iterate(yield)
}
