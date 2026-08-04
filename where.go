package linq

// Where filters a collection of values based on a predicate.
func (q Query[T]) Where(predicate func(T) bool) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			q.Iterate(func(item T) bool {
				if predicate(item) {
					return yield(item)
				}
				return true
			})
		},
	}
}

// WhereIndexed filters a collection of values based on a predicate. Each
// element's index is used in the logic of the predicate function.
//
// The first argument represents the zero-based index of the element within
// the collection. The second argument of predicate represents the element to test.
func (q Query[T]) WhereIndexed(predicate func(int, T) bool) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			index := 0
			q.Iterate(func(item T) bool {
				shouldYield := predicate(index, item)
				index++

				if shouldYield {
					return yield(item)
				}

				return true
			})
		},
	}
}
