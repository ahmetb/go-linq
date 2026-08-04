package linq

// Reverse inverts the order of the elements in a collection.
//
// Unlike OrderBy, this sorting method does not consider the actual values
// themselves in determining the order. Rather, it just returns the elements in
// the reverse order from which they are produced by the underlying source.
func (q Query[T]) Reverse() Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			var items []T
			for item := range q.Iterate {
				items = append(items, item)
			}

			for i := len(items) - 1; i >= 0; i-- {
				if !yield(items[i]) {
					return
				}
			}
		},
	}
}
