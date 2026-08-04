package linq

func (q Query[T]) ForEach(action func(T)) {
	for value := range q.iterate {
		action(value)
	}
}
