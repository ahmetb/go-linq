package linq

// All reports whether every value satisfies predicate; it is true for an empty query.
func (q Query[T]) All(predicate func(T) bool) bool {
	for value := range q.iterate {
		if !predicate(value) {
			return false
		}
	}

	return true
}
