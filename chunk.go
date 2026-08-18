package linq

import "slices"

// Chunk splits the elements of a collection into slices of at most size
// elements. Every returned slice has its own backing storage. Chunk panics if
// size is less than one.
func Chunk[T any](q Query[T], size int) Query[[]T] {
	if size < 1 {
		panic("linq: chunk size must be positive")
	}

	resultSize := 0
	if q.size > 0 {
		resultSize = (q.size-1)/size + 1
	}

	return Query[[]T]{
		Iterate: func(yield func([]T) bool) {
			var chunk []T
			q.Iterate(func(item T) bool {
				if chunk == nil {
					chunk = make([]T, 0, size)
				}
				chunk = append(chunk, item)
				if len(chunk) < size {
					return true
				}

				current := chunk
				chunk = nil
				return yield(current)
			})

			if len(chunk) > 0 {
				yield(slices.Clip(chunk))
			}
		},
		size: resultSize,
	}
}
