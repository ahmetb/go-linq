package linq

// Count returns the number of values in the query.
func (q Query[T]) Count() int {
	count := 0
	q.iterate(func(T) bool {
		count++
		return true
	})
	return count
}
