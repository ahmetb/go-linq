package linq

// TakeWhileIndexed is TakeWhile with each zero-based source index.
func (q Query[T]) TakeWhileIndexed(predicate func(int, T) bool) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		index := 0
		q.iterate(func(value T) bool {
			take := predicate(index, value)
			index++
			if !take {
				return false
			}

			return yield(value)
		})
	}}
}
