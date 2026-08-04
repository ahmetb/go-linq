package linq

// Single returns the sole value, or the zero value of T unless exactly one exists.
func (q Query[T]) Single() T {
	var result T
	found := false
	q.iterate(func(value T) bool {
		if found {
			var zero T
			result = zero
			return false
		}

		result = value
		found = true
		return true
	})
	return result
}
