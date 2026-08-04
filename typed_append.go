package linq

func (q Query[T]) Append(value T) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
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
