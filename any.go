package linq

func (q Query[T]) Any() bool {
	found := false
	q.iterate(func(T) bool {
		found = true
		return false
	})
	return found
}
