package linq

// Last returns the final value, or the zero value of T when empty.
func (q Query[T]) Last() T {
	var result T
	q.iterate(func(value T) bool {
		result = value
		return true
	})
	return result
}
