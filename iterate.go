package linq

// Iterate calls yield for each value until the source ends or yield returns false.
func (q Query[T]) Iterate(yield func(T) bool) {
	q.iterate(yield)
}
