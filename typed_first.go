package linq

func (q query[T]) First() T {
	var result T
	q.iterate(func(value T) bool {
		result = value
		return false
	})
	return result
}
