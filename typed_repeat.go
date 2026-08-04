package linq

func repeatValue[T any](value T, count int) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		for range count {
			if !yield(value) {
				return
			}
		}
	}}
}
