package linq

func (q query[T]) Append(value T) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		continuing := true
		q.iterate(func(current T) bool {
			continuing = yield(current)
			return continuing
		})
		if continuing {
			yield(value)
		}
	}}
}
