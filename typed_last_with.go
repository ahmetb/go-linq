package linq

func (q query[T]) LastWith(predicate func(T) bool) T {
	var result T
	for value := range q.iterate {
		if predicate(value) {
			result = value
		}
	}

	return result
}
