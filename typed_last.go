package linq

func (q query[T]) Last() T {
	var result T
	q.iterate(func(value T) bool {
		result = value
		return true
	})
	return result
}
