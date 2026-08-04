package linq

func repeatValue[T any](value T, count int) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		for range count {
			if !yield(value) {
				return
			}
		}
	}}
}
