package linq

// Take returns a specified number of contiguous elements from the start of a
// collection. It stops pulling from the source as soon as count elements have
// been yielded.
func (q Query[T]) Take(count int) Query[T] {
	if count <= 0 {
		return Query[T]{Iterate: func(func(T) bool) {}}
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
