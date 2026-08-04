package linq

func (q query[T]) Any() bool {
	found := false
	q.iterate(func(T) bool {
		found = true
		return false
	})
	return found
}
