package linq

// IntersectByT is the typed version of IntersectBy.
//
// NOTE: IntersectBy method has better performance than IntersectByT
//
// selectorFn is of a type "func(TSource) TSource"
//
func (q Query) IntersectByT(q2 Query, selectorFn interface{}) Query {

	selectorGenericFunc, err := newGenericFunc(
		"IntersectByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.IntersectBy(q2, selectorFunc)
}
