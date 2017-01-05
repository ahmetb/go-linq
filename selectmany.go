package linq

// SelectMany projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection.
func (q Query) SelectMany(selector func(interface{}) Query) Query {
	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			var inner interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if inner == nil {
						inner, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(inner).Iterate()
					}

					item, ok = innernext()
					if !ok {
						inner = nil
					}
				}

				return
			}
		},
	}
}

// SelectManyT is the typed version of SelectMany.
//
//   - selectorFn is of type "func(TSource)Query"
//
// NOTE: SelectMany has better performance than SelectManyT.
func (q Query) SelectManyT(selectorFn interface{}) Query {

	selectManyGenericFunc, err := newGenericFunc(
		"SelectManyT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(Query))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(inner interface{}) Query {
		return selectManyGenericFunc.Call(inner).(Query)
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
func (q Query) SelectManyIndexed(selector func(int, interface{}) Query) Query {
	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			index := 0
			var inner interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if inner == nil {
						inner, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(index, inner).Iterate()
						index++
					}

					item, ok = innernext()
					if !ok {
						inner = nil
					}
				}

				return
			}
		},
	}
}

// SelectManyIndexedT is the typed version of SelectManyIndexed.
//
//   - selectorFn is of type "func(int,TSource)Query"
//
// NOTE: SelectManyIndexed has better performance than SelectManyIndexedT.
func (q Query) SelectManyIndexedT(selectorFn interface{}) Query {

	selectManyIndexedGenericFunc, err := newGenericFunc(
		"SelectManyIndexedT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(Query))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(index int, inner interface{}) Query {
		return selectManyIndexedGenericFunc.Call(index, inner).(Query)
	}

	return q.SelectManyIndexed(selectorFunc)
}

// SelectManyBy projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection, and invokes a result
// selector function on each element therein.
func (q Query) SelectManyBy(selector func(interface{}) Query,
	resultSelector func(interface{}, interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			var outer interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if outer == nil {
						outer, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(outer).Iterate()
					}

					item, ok = innernext()
					if !ok {
						outer = nil
					}
				}

				item = resultSelector(item, outer)
				return
			}
		},
	}
}

// SelectManyByT is the typed version of SelectManyBy.
//
//   - selectorFn is of type "func(TSource)Query"
//   - resultSelectorFn is of type "func(TSource,TCollection)TResult"
//
// NOTE: SelectManyBy has better performance than SelectManyByT.
func (q Query) SelectManyByT(selectorFn interface{},
	resultSelectorFn interface{}) Query {

	selectorGenericFunc, err := newGenericFunc(
		"SelectManyByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(Query))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(outer interface{}) Query {
		return selectorGenericFunc.Call(outer).(Query)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"SelectManyByT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer interface{}, item interface{}) interface{} {
		return resultSelectorGenericFunc.Call(outer, item)
	}

	return q.SelectManyBy(selectorFunc, resultSelectorFunc)
}

// SelectManyByIndexed projects each element of a collection to a Query,
// iterates and flattens the resulting collection into one collection, and
// invokes a result selector function on each element therein. The index of each
// source element is used in the intermediate projected form of that element.
func (q Query) SelectManyByIndexed(selector func(int, interface{}) Query,
	resultSelector func(interface{}, interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			index := 0
			var outer interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if outer == nil {
						outer, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(index, outer).Iterate()
						index++
					}

					item, ok = innernext()
					if !ok {
						outer = nil
					}
				}

				item = resultSelector(item, outer)
				return
			}
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
func (q Query) SelectManyByIndexedT(selectorFn interface{},
	resultSelectorFn interface{}) Query {
	selectorGenericFunc, err := newGenericFunc(
		"SelectManyByIndexedT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(Query))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(index int, outer interface{}) Query {
		return selectorGenericFunc.Call(index, outer).(Query)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"SelectManyByIndexedT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer interface{}, item interface{}) interface{} {
		return resultSelectorGenericFunc.Call(outer, item)
	}

	return q.SelectManyByIndexed(selectorFunc, resultSelectorFunc)
}
