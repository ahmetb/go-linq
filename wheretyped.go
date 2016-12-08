package linq

// WhereT is the typed version of Where.
//
// NOTE: Where method has better performance than WhereT
//
// predicateFn is of a type "func(TSource)bool"
//
func (q Query) WhereT(predicateFn interface{}) Query {

	predicateGenericFunc, err := newGenericFunc(
		"WhereT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.Where(predicateFunc)
}

// WhereIndexedT is the typed version of WhereIndexed.
//
// NOTE: WhereIndexed method has better performance than WhereIndexedT
//
// predicateFn is of a type "func(int,TSource)bool"
//
func (q Query) WhereIndexedT(predicateFn interface{}) Query {

	predicateGenericFunc, err := newGenericFunc(
		"WhereIndexedT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(index int, item interface{}) bool {
		return predicateGenericFunc.Call(index, item).(bool)
	}

	return q.WhereIndexed(predicateFunc)
}
