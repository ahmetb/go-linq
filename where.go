package linq

func (q Query[T]) Where(predicate func(T) bool) Query[T] {
	return Query[T]{
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

func (q Query[T]) WhereIndexed(predicate func(int, T) bool) Query[T] {
	return Query[T]{
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
