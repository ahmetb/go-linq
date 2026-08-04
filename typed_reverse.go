package linq

import "slices"

func (q query[T]) Reverse() query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		values := q.toSlice()
		slices.Reverse(values)
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}}
}
