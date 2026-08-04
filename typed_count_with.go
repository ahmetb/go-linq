package linq

func (q query[T]) CountWith(predicate func(T) bool) int {
	count := 0
	for value := range q.iterate {
		if predicate(value) {
			count++
		}
	}

	return count
}
