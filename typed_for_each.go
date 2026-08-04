package linq

func (q query[T]) ForEach(action func(T)) {
	for value := range q.iterate {
		action(value)
	}
}
