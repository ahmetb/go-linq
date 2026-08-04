package linq

func (q query[T]) IndexOf(predicate func(T) bool) int {
	result := -1
	index := 0
	q.iterate(func(value T) bool {
		if predicate(value) {
			result = index
			return false
		}

		index++
		return true
	})
	return result
}
