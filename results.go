package linq

// Results materializes the query into a new typed slice.
func (q Query[T]) Results() []T {
	return q.toSlice()
}
