package linq

func (q Query[T]) First() T {
	var result T
	q.iterate(func(value T) bool {
		result = value
		return false
	})
	return result
}
