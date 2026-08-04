package linq

// SingleWith returns the sole matching value, or the zero value of T otherwise.
func (q Query[T]) SingleWith(predicate func(T) bool) T {
	var result T
	found := false
	for value := range q.iterate {
		if !predicate(value) {
			continue
		}
		if found {
			var zero T
			return zero
		}

		result = value
		found = true
	}

	if found {
		return result
	}

	var zero T
	return zero
}
