package linq

// JoinT is the typed version of Join.
//
// NOTE: Join method has better performance than JoinT
//
// outerKeySelectorFn is of a type "func(TOuter) TKey"
//
// innerKeySelectorFn is of a type "func(TInner) TKey"
//
// resultSelectorFn is of a type "func(TOuter,TInner) TResult"
//
func (q Query) JoinT(inner Query,
	outerKeySelectorFn interface{},
	innerKeySelectorFn interface{},
	resultSelectorFn interface{},
) Query {

	outerKeySelectorGenericFunc, err := newGenericFunc(
		"JoinT", "outerKeySelectorFn", outerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	outerKeySelectorFunc := func(item interface{}) interface{} {
		return outerKeySelectorGenericFunc.Call(item)
	}

	innerKeySelectorFuncGenericFunc, err := newGenericFunc(
		"JoinT", "innerKeySelectorFn",
		innerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	innerKeySelectorFunc := func(item interface{}) interface{} {
		return innerKeySelectorFuncGenericFunc.Call(item)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"JoinT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer interface{}, inner interface{}) interface{} {
		return resultSelectorGenericFunc.Call(outer, inner)
	}

	return q.Join(inner, outerKeySelectorFunc, innerKeySelectorFunc, resultSelectorFunc)
}
