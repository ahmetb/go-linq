package linq

// ToChannel synchronously sends every value to result and then closes it.
func (q Query[T]) ToChannel(result chan<- T) {
	defer close(result)

	q.iterate(func(value T) bool {
		result <- value
		return true
	})
}
