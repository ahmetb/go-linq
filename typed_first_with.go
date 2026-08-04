package linq

func (q query[T]) FirstWith(predicate func(T) bool) T {
	for value := range q.iterate {
		if predicate(value) {
			return value
		}
	}

	var zero T
	return zero
}
