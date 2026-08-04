package linq

// SelectMany projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection.
func (q Query[T]) SelectMany[TResult any](selector func(T) Query[TResult]) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			q.Iterate(func(outerItem T) bool {
				keepGoing := true

				innerQuery := selector(outerItem)
				innerQuery.Iterate(func(innerItem TResult) bool {
					if !yield(innerItem) {
						keepGoing = false
						return false
					}
					return true
				})

				return keepGoing
			})
		},
	}
}

// SelectManyIndexed projects each element of a collection to a Query, iterates
// and flattens the resulting collection into one collection.
//
// The first argument to selector represents the zero-based index of that
// element in the source collection. This can be useful if the elements are in a
// known order and you want to do something with an element at a particular
// index, for example. It can also be useful if you want to retrieve the index
// of one or more elements. The second argument to selector represents the
// element to process.
func (q Query[T]) SelectManyIndexed[TResult any](selector func(index int, outer T) Query[TResult]) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			index := 0
			q.Iterate(func(outerItem T) bool {
				keepGoing := true

				innerQuery := selector(index, outerItem)
				index++
				innerQuery.Iterate(func(innerItem TResult) bool {
					if !yield(innerItem) {
						keepGoing = false
						return false
					}
					return true
				})

				return keepGoing
			})
		},
	}
}

// SelectManyBy projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection, and invokes a result
// selector function on each element therein.
func (q Query[T]) SelectManyBy[TCollection, TResult any](
	selector func(outer T) Query[TCollection],
	resultSelector func(inner TCollection, outer T) TResult,
) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			q.Iterate(func(outerItem T) bool {
				keepGoing := true
				innerQuery := selector(outerItem)

				innerQuery.Iterate(func(innerItem TCollection) bool {
					result := resultSelector(innerItem, outerItem)

					if !yield(result) {
						keepGoing = false
						return false
					}
					return true
				})

				return keepGoing
			})
		},
	}
}

// SelectManyByIndexed projects each element of a collection to a Query,
// iterates and flattens the resulting collection into one collection, and
// invokes a result selector function on each element therein. The index of each
// source element is used in the intermediate projected form of that element.
func (q Query[T]) SelectManyByIndexed[TCollection, TResult any](
	selector func(index int, outer T) Query[TCollection],
	resultSelector func(inner TCollection, outer T) TResult,
) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			index := 0
			q.Iterate(func(outerItem T) bool {
				innerQuery := selector(index, outerItem)
				index++

				keepGoing := true
				innerQuery.Iterate(func(innerItem TCollection) bool {
					result := resultSelector(innerItem, outerItem)

					if !yield(result) {
						keepGoing = false
						return false
					}
					return true
				})

				return keepGoing
			})
		},
	}
}
