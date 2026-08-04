package linq

// AnyWith reports whether at least one value satisfies predicate.
func (q Query[T]) AnyWith(predicate func(T) bool) bool {
	for value := range q.iterate {
		if predicate(value) {
			return true
		}
	}

	return false
}
