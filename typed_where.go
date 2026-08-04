package linq

func (q query[T]) Where(predicate func(T) bool) query[T] {
	return query[T]{
		iterate: func(yield func(T) bool) {
			q.iterate(func(value T) bool {
				if predicate(value) {
					return yield(value)
				}
				return true
			})
		},
	}
}

func (q query[T]) WhereIndexed(predicate func(int, T) bool) query[T] {
	return query[T]{
		iterate: func(yield func(T) bool) {
			index := 0
			q.iterate(func(value T) bool {
				accepted := predicate(index, value)
				index++
				if accepted {
					return yield(value)
				}
				return true
			})
		},
	}
}
