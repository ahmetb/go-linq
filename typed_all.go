package linq

func (q query[T]) All(predicate func(T) bool) bool {
	for value := range q.iterate {
		if !predicate(value) {
			return false
		}
	}

	return true
}
