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

// chunkBuffer holds the current chunk and the next buffer's initial capacity.
type chunkBuffer[T any] struct {
	items    []T
	capacity int
}

// chunked yields independent chunks without constructing Query[[]T], which
// would cause a recursive instantiation inside ChunkBy.
func chunked[T any](q Query[T], size int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		// Unknown sources start small. A hint bounds the initial reservation,
		// but never limits how many elements we actually read.
		buffer := chunkBuffer[T]{capacity: min(size, 4)}
		if q.size > 0 {
			buffer.capacity = min(size, q.size)
		}
		q.Iterate(func(item T) bool {
			if buffer.items == nil {
				buffer.items = make([]T, 0, buffer.capacity)
			}
			buffer.items = append(buffer.items, item)
			if len(buffer.items) < size {
				return true
			}

			current := buffer.items
			buffer.items = nil
			buffer.capacity = size // a full chunk established that this size is useful
			if cap(current) > len(current) {
				current = slices.Clone(current)
			}
			return yield(slices.Clip(current))
		})

		if chunk := buffer.items; len(chunk) > 0 {
			if len(chunk) < cap(chunk) {
				// Clip alone retains the entire backing allocation.
				chunk = slices.Clone(chunk)
			}
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
