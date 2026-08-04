package linq

// SelectMany projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection.
func (q legacyQuery) SelectMany(selector func(any) legacyQuery) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			q.Iterate(func(outerItem any) bool {
				keepGoing := true

				innerQuery := selector(outerItem)
				innerQuery.Iterate(func(innerItem any) bool {
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

// SelectManyT is the typed version of SelectMany.
//
//   - selectorFn is of type "func(TSource)Query"
//
// NOTE: SelectMany has better performance than SelectManyT.
func (q legacyQuery) SelectManyT(selectorFn any) legacyQuery {

	selectManyGenericFunc, err := newGenericFunc(
		"SelectManyT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(legacyQuery))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(inner any) legacyQuery {
		return selectManyGenericFunc.Call(inner).(legacyQuery)
	}
	return q.SelectMany(selectorFunc)

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
func (q legacyQuery) SelectManyIndexed(selector func(index int, outer any) legacyQuery) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			index := 0
			q.Iterate(func(outerItem any) bool {
				keepGoing := true

				innerQuery := selector(index, outerItem)
				index++
				innerQuery.Iterate(func(innerItem any) bool {
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

// SelectManyIndexedT is the typed version of SelectManyIndexed.
//
//   - selectorFn is of type "func(int,TSource)Query"
//
// NOTE: SelectManyIndexed has better performance than SelectManyIndexedT.
func (q legacyQuery) SelectManyIndexedT(selectorFn any) legacyQuery {

	selectManyIndexedGenericFunc, err := newGenericFunc(
		"SelectManyIndexedT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(legacyQuery))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(index int, inner any) legacyQuery {
		return selectManyIndexedGenericFunc.Call(index, inner).(legacyQuery)
	}

	return q.SelectManyIndexed(selectorFunc)
}

// SelectManyBy projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection, and invokes a result
// selector function on each element therein.
func (q legacyQuery) SelectManyBy(
	selector func(outer any) legacyQuery, resultSelector func(inner, outer any) any,
) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			q.Iterate(func(outerItem any) bool {
				keepGoing := true
				innerQuery := selector(outerItem)

				innerQuery.Iterate(func(innerItem any) bool {
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

// SelectManyByT is the typed version of SelectManyBy.
//
//   - selectorFn is of type "func(TSource)Query"
//   - resultSelectorFn is of type "func(TSource,TCollection)TResult"
//
// NOTE: SelectManyBy has better performance than SelectManyByT.
func (q legacyQuery) SelectManyByT(selectorFn any,
	resultSelectorFn any) legacyQuery {

	selectorGenericFunc, err := newGenericFunc(
		"SelectManyByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(legacyQuery))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(outer any) legacyQuery {
		return selectorGenericFunc.Call(outer).(legacyQuery)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"SelectManyByT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer any, item any) any {
		return resultSelectorGenericFunc.Call(outer, item)
	}

	return q.SelectManyBy(selectorFunc, resultSelectorFunc)
}

// SelectManyByIndexed projects each element of a collection to a Query,
// iterates and flattens the resulting collection into one collection, and
// invokes a result selector function on each element therein. The index of each
// source element is used in the intermediate projected form of that element.
func (q legacyQuery) SelectManyByIndexed(
	selector func(index int, outer any) legacyQuery, resultSelector func(inner, outer any) any,
) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			index := 0
			q.Iterate(func(outerItem any) bool {
				innerQuery := selector(index, outerItem)
				index++

				keepGoing := true
				innerQuery.Iterate(func(innerItem any) bool {
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

// SelectManyByIndexedT is the typed version of SelectManyByIndexed.
//
//   - selectorFn is of type "func(int,TSource)Query"
//   - resultSelectorFn is of type "func(TSource,TCollection)TResult"
//
// NOTE: SelectManyByIndexed has better performance than
// SelectManyByIndexedT.
func (q legacyQuery) SelectManyByIndexedT(selectorFn any,
	resultSelectorFn any) legacyQuery {
	selectorGenericFunc, err := newGenericFunc(
		"SelectManyByIndexedT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(legacyQuery))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(index int, outer any) legacyQuery {
		return selectorGenericFunc.Call(index, outer).(legacyQuery)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"SelectManyByIndexedT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer any, item any) any {
		return resultSelectorGenericFunc.Call(outer, item)
	}

	return q.SelectManyByIndexed(selectorFunc, resultSelectorFunc)
}
