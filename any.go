package linq

// Any reports whether the query contains at least one value.
func (q Query[T]) Any() bool {
	found := false
	q.iterate(func(T) bool {
		found = true
		return false
	})
	return found
}
