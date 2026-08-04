package linq

// First returns the first value, or the zero value of T when empty.
func (q Query[T]) First() T {
	var result T
	q.iterate(func(value T) bool {
		result = value
		return false
	})
	return result
}
