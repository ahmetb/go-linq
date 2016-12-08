package linq

// GroupByT is the typed version of GroupBy.
//
// NOTE: GroupBy method has better performance than GroupByT
//
// keySelectorFn is of a type "func(TSource) TKey"
//
// elementSelectorFn is of a type "func(TSource) TElement"
//
func (q Query) GroupByT(keySelectorFn interface{}, elementSelectorFn interface{}) Query {

	keySelectorGenericFunc, err := newGenericFunc(
		"GroupByT", "keySelectorFn", keySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	keySelectorFunc := func(item interface{}) interface{} {
		return keySelectorGenericFunc.Call(item)
	}

	elementSelectorGenericFunc, err := newGenericFunc(
		"GroupByT", "elementSelectorFn", elementSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	elementSelectorFunc := func(item interface{}) interface{} {
		return elementSelectorGenericFunc.Call(item)

	}

	return q.GroupBy(keySelectorFunc, elementSelectorFunc)
}
