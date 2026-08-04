package linq

// ForEachIndexed invokes action with each zero-based source index and value.
func (q Query[T]) ForEachIndexed(action func(int, T)) {
	index := 0
	for value := range q.iterate {
		action(index, value)
		index++
	}
}
