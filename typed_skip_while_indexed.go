package linq

func (q Query[T]) SkipWhileIndexed(predicate func(int, T) bool) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
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
