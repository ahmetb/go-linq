package linq

func (q query[T]) TakeWhileIndexed(predicate func(int, T) bool) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
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
