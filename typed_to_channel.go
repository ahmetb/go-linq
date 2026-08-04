package linq

func (q query[T]) ToChannel(result chan<- T) {
	defer close(result)

	q.iterate(func(value T) bool {
		result <- value
		return true
	})
}
