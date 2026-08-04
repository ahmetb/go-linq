package linq

// Skip bypasses a specified number of elements in a collection and then returns
// the remaining elements.
func (q Query[T]) Skip(count int) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			n := count
			q.Iterate(func(item T) bool {
				if n > 0 {
					n--
					return true
				}
				return yield(item)
			})
		},
	}
}

// SkipWhile bypasses elements in a collection as long as a specified condition
// is true and then returns the remaining elements.
//
// This method tests each element by using predicate and skips the element if
// the result is true. After the predicate function returns false for an
// element, that element and the remaining elements in source are returned and
// there are no more invocations of predicate.
func (q Query[T]) SkipWhile(predicate func(T) bool) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			skipping := true
			q.Iterate(func(item T) bool {
				if skipping {
					if predicate(item) {
						return true
					}
					skipping = false
				}

				return yield(item)
			})
		},
	}
}

// SkipWhileIndexed bypasses elements in a collection as long as a specified
// condition is true and then returns the remaining elements. The element's
// index is used in the logic of the predicate function.
//
// This method tests each element by using predicate and skips the element if
// the result is true. After the predicate function returns false for an
// element, that element and the remaining elements in source are returned and
// there are no more invocations of predicate.
func (q Query[T]) SkipWhileIndexed(predicate func(int, T) bool) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			skipping := true
			index := 0
			q.Iterate(func(item T) bool {
				if skipping {
					if predicate(index, item) {
						index++
						return true
					}
					skipping = false
				}

				return yield(item)
			})
		},
	}
}
