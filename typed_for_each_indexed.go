package linq

func (q query[T]) ForEachIndexed(action func(int, T)) {
	index := 0
	for value := range q.iterate {
		action(index, value)
		index++
	}
}
