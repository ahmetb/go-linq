package linq

func fromChannel[T any](source <-chan T) Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		for value := range source {
			if !yield(value) {
				return
			}
		}
	}}
}
