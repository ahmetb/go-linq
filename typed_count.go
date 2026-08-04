package linq

func (q query[T]) Count() int {
	count := 0
	q.iterate(func(T) bool {
		count++
		return true
	})
	return count
}
