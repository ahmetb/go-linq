package linq

// SelectManyT is the typed version of SelectMany.
//
// NOTE: SelectMany method has better performance than SelectManyT
//
// selectorFn is of a type "func(TSource)Query"
//
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

// SelectManyIndexedT is the typed version of SelectManyIndexed.
//
// NOTE: SelectManyIndexed method has better performance than SelectManyIndexedT
//
// selectorFn is of a type "func(int,TSource)Query"
//
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

// SelectManyByT is the typed version of SelectManyBy.
//
// NOTE: SelectManyBy method has better performance than SelectManyByT
//
// selectorFn is of a type "func(TSource)Query"
//
// resultSelectorFn is of a type "func(TSource,TCollection)TResult"
//
func (q Query) SelectManyByT(selectorFn interface{}, resultSelectorFn interface{}) Query {

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

// SelectManyByIndexedT is the typed version of SelectManyByIndexed.
//
// NOTE: SelectManyByIndexed method has better performance than SelectManyByIndexedT
//
// selectorFn is of a type "func(int,TSource)Query"
//
// resultSelectorFn is of a type "func(TSource,TCollection)TResult"
//
func (q Query) SelectManyByIndexedT(selectorFn interface{}, resultSelectorFn interface{}) Query {

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
