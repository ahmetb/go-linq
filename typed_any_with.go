package linq

func (q query[T]) AnyWith(predicate func(T) bool) bool {
	for value := range q.iterate {
		if predicate(value) {
			return true
		}
	}

	return false
}
