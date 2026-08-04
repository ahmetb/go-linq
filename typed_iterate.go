package linq

func (q query[T]) Iterate(yield func(T) bool) {
	q.iterate(yield)
}
