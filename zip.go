package linq

import "iter"

// Zip applies a specified function to the corresponding elements of two
// collections, producing a collection of the results.
//
// The method steps through the two input collections, applying function
// resultSelector to corresponding elements of the two collections. The method
// returns a collection of the values that are returned by resultSelector. If
// the input collections do not have the same number of elements, the method
// combines elements until it reaches the end of one of the collections. For
// example, if one collection has three elements and the other one has four, the
// result collection has only three elements.
func (q Query[T]) Zip[TSecond, TResult any](q2 Query[TSecond],
	resultSelector func(T, TSecond) TResult) Query[TResult] {

	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			next2, stop2 := iter.Pull(q2.Iterate)
			defer stop2()

			q.Iterate(func(item T) bool {
				item2, ok2 := next2()
				if !ok2 {
					return false
				}

				return yield(resultSelector(item, item2))
			})
		},
	}
}
