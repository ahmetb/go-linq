package linq

// SelectT is the typed version of Select.
//
// NOTE: Select method has better performance than SelectT
//
// selectorFn is of a type "func(TSource)TResult"
//
func (q Query) SelectT(selectorFn interface{}) Query {

	selectGenericFunc, err := newGenericFunc(
		"SelectT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectGenericFunc.Call(item)
	}

	return q.Select(selectorFunc)
}

// SelectIndexedT is the typed version of SelectIndexed.
//
// NOTE: SelectIndexed method has better performance than SelectIndexedT
//
// selectorFn is of a type "func(int,TSource)TResult"
//
func (q Query) SelectIndexedT(selectorFn interface{}) Query {

	selectGenericFunc, err := newGenericFunc(
		"SelectIndexedT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(index int, item interface{}) interface{} {
		return selectGenericFunc.Call(index, item)
	}

	return q.SelectIndexed(selectorFunc)
}
