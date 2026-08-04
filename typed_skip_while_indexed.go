package linq

func (q query[T]) SkipWhileIndexed(predicate func(int, T) bool) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		index := 0
		skipping := true
		q.iterate(func(value T) bool {
			if skipping {
				skip := predicate(index, value)
				index++
				if skip {
					return true
				}
				skipping = false
			}

			return yield(value)
		})
	}}
}
