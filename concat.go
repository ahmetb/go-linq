package linq

// Append inserts an item to the end of a collection, so it becomes the last
// item.
func (q Query[T]) Append(item T) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			stopped := false

			q.Iterate(func(originalItem T) bool {
				if !yield(originalItem) {
					stopped = true
					return false
				}
				return true
			})

			if !stopped {
				yield(item)
			}
		},
	}
}

// Concat concatenates two collections.
//
// The Concat method differs from the Union method because the Concat method
// returns all the original elements in the input sequences. The Union method
// returns only unique elements.
func (q Query[T]) Concat(q2 Query[T]) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			stopped := false

			q.Iterate(func(item T) bool {
				if !yield(item) {
					stopped = true
					return false
				}
				return true
			})

			if !stopped {
				q2.Iterate(yield)
			}
		},
	}
}

// Prepend inserts an item to the beginning of a collection, so it becomes the
// first item.
func (q Query[T]) Prepend(item T) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			if !yield(item) {
				return
			}

			q.Iterate(yield)
		},
	}
}
