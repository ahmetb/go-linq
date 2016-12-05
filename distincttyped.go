package linq

// DistinctByT is the typed version of DistinctBy.
//
// NOTE: DistinctBy method has better performance than DistinctByT
//
// selectorFn is of type "func(TSource) TSource".
//
func (q Query) DistinctByT(selectorFn interface{}) Query {
	selectorFunc, ok := selectorFn.(func(interface{}) interface{})
	if !ok {
		selectorGenericFunc, err := newGenericFunc(
			"DistinctBy", "selectorFn", selectorFn,
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
		)
		if err != nil {
			panic(err)
		}

		selectorFunc = func(item interface{}) interface{} {
			return selectorGenericFunc.Call(item)
		}
	}
	return q.DistinctBy(selectorFunc)
}
