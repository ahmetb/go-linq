package linq

// CountWith returns the number of values satisfying predicate.
func (q Query[T]) CountWith(predicate func(T) bool) int {
	count := 0
	for value := range q.iterate {
		if predicate(value) {
			count++
		}
	}

	return count
}
