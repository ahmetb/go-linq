package linq

import "slices"

func (q Query[T]) Reverse() Query[T] {
	return Query[T]{iterate: func(yield func(T) bool) {
		values := q.toSlice()
		slices.Reverse(values)
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}}
}
