package linq

// ExceptByT is the typed version of ExceptByT.
//
// NOTE: ExceptBy method has better performance than ExceptByT
//
// selectorFn is of a type "func(TSource) TSource"
//
func (q Query) ExceptByT(q2 Query, selectorFn interface{}) Query {

	selectorGenericFunc, err := newGenericFunc(
		"ExceptByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.ExceptBy(q2, selectorFunc)
}
