package linq

func fromChannel[T any](source <-chan T) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		for value := range source {
			if !yield(value) {
				return
			}
		}
	}}
}
