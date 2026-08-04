package linq

// FirstWith returns the first matching value, or the zero value of T.
func (q Query[T]) FirstWith(predicate func(T) bool) T {
	for value := range q.iterate {
		if predicate(value) {
			return value
		}
	}

	var zero T
	return zero
}
