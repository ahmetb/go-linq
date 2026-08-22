package linq

import (
	"iter"
	"slices"
)

// Chunk splits the elements of a collection into slices of at most size
// elements. Every returned slice has its own backing storage. Chunk panics if
// size is less than one.
//
// Chunk is a package-level function because a Query[T].Chunk() Query[[]T]
// method would be a recursive instantiation. To split a collection without
// leaving a chain, use the ChunkBy method.
func Chunk[T any](q Query[T], size int) Query[[]T] {
	if size < 1 {
		panic("linq: chunk size must be positive")
	}

	return Query[[]T]{
		Iterate: chunked(q, size),
		size:    chunkCount(q.size, size),
	}
}

// ChunkBy splits the elements of a collection into slices of at most size
// elements and projects each slice with selector. It is Chunk followed by
// Select, available as a method so that it composes inside a chain. ChunkBy
// panics if size is less than one.
//
// Each slice passed to selector has its own backing storage and is not reused
// afterwards, so selector may retain it.
//
// ChunkBy is a generic method: the result type TResult is inferred from the
// selector function.
func (q Query[T]) ChunkBy[TResult any](size int, selector func([]T) TResult) Query[TResult] {
	if size < 1 {
		panic("linq: chunk size must be positive")
	}

	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			chunked(q, size)(func(chunk []T) bool {
				return yield(selector(chunk))
			})
		},
		size: chunkCount(q.size, size),
	}
}

// chunked yields independent chunks without constructing Query[[]T], which
// would cause a recursive instantiation inside ChunkBy.
func chunked[T any](q Query[T], size int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
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
	}
}

func chunkCount(length, size int) int {
	if length <= 0 {
		return 0
	}
	return (length-1)/size + 1
}
