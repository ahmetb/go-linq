package linq

import "slices"

// Take returns a specified number of contiguous elements from the start of a
// collection. It stops pulling from the source as soon as count elements have
// been yielded.
func (q Query[T]) Take(count int) Query[T] {
	if count <= 0 {
		return emptyQuery[T]()
	}

	return Query[T]{
		Iterate: func(yield func(T) bool) {
			n := count
			q.Iterate(func(item T) bool {
				if !yield(item) {
					return false
				}
				n--
				return n > 0
			})
		},
	}
}

// lastItems consumes q and returns its last count elements in their original
// order together with the total number of elements consumed. Its storage is
// bounded by count.
func lastItems[T any](q Query[T], count int) ([]T, int) {
	if count <= 0 {
		return nil, 0
	}

	capacity := q.size
	if capacity == 0 || capacity > count {
		capacity = count
	}
	items := make([]T, 0, capacity)
	head := 0
	total := 0
	q.Iterate(func(item T) bool {
		total++
		if len(items) < count {
			items = append(items, item)
			return true
		}

		items[head] = item
		head++
		if head == len(items) {
			head = 0
		}
		return true
	})

	if head == 0 {
		return items, total
	}

	slices.Reverse(items[:head])
	slices.Reverse(items[head:])
	slices.Reverse(items)
	return items, total
}

// TakeLast returns the last count elements of a collection in their original
// order. It consumes the source before yielding its first element.
func (q Query[T]) TakeLast(count int) Query[T] {
	if count <= 0 {
		return emptyQuery[T]()
	}

	return Query[T]{
		Iterate: func(yield func(T) bool) {
			items, _ := lastItems(q, count)
			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		},
		size: min(count, q.size),
	}
}

// TakeRange returns the contiguous elements between start (inclusive) and end
// (exclusive). Out-of-range bounds are clipped; a reversed range is empty. A
// start counted from the end consumes the entire source before yielding its
// first element.
func (q Query[T]) TakeRange(start, end Position) Query[T] {
	resultSize := max(end.offset(q.size)-start.offset(q.size), 0)

	if !start.fromEnd {
		if !end.fromEnd && start.value >= end.value {
			return emptyQuery[T]()
		}

		result := q.Skip(start.value)
		if end.fromEnd {
			result = result.SkipLast(end.value)
		} else {
			result = result.Take(end.value - start.value)
		}
		result.size = resultSize
		return result
	}

	if start.value == 0 ||
		(end.fromEnd && end.value >= start.value) ||
		(!end.fromEnd && end.value == 0) {
		return emptyQuery[T]()
	}

	return Query[T]{
		Iterate: func(yield func(T) bool) {
			items, total := lastItems(q, start.value)
			length := max(end.offset(total)-start.offset(total), 0)

			for _, item := range items[:length] {
				if !yield(item) {
					return
				}
			}
		},
		size: resultSize,
	}
}

// TakeWhile returns elements from a collection as long as a specified condition
// is true and then skips the remaining elements.
func (q Query[T]) TakeWhile(predicate func(T) bool) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			q.Iterate(func(item T) bool {
				if predicate(item) {
					return yield(item)
				}
				return false
			})
		},
	}
}

// TakeWhileIndexed returns elements from a collection as long as a specified
// condition is true. The element's index is used in the logic of the predicate
// function. The first argument of predicate represents the zero-based index of
// the element within the collection. The second argument represents the element to
// test.
func (q Query[T]) TakeWhileIndexed(predicate func(int, T) bool) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			index := 0
			q.Iterate(func(item T) bool {
				if predicate(index, item) {
					index++
					return yield(item)
				}
				return false
			})
		},
	}
}
