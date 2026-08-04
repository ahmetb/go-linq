package linq

// LastWith returns the final matching value, or the zero value of T.
func (q Query[T]) LastWith(predicate func(T) bool) T {
	var result T
	for value := range q.iterate {
		if predicate(value) {
			result = value
		}
	}

	return result
}
