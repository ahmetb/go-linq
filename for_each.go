package linq

// ForEach invokes action once for every value in source order.
func (q Query[T]) ForEach(action func(T)) {
	for value := range q.iterate {
		action(value)
	}
}
