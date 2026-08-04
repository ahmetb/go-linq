package linq

func (q query[T]) Single() T {
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
