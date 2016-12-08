package linq

// TakeWhileT is the typed version of TakeWhile.
//
// NOTE: TakeWhile method has better performance than TakeWhileT
//
// predicateFn is of a type "func(TSource)bool"
//
func (q Query) TakeWhileT(predicateFn interface{}) Query {

	predicateGenericFunc, err := newGenericFunc(
		"TakeWhileT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.TakeWhile(predicateFunc)
}

// TakeWhileIndexedT is the typed version of TakeWhileIndexed.
//
// NOTE: TakeWhileIndexed method has better performance than TakeWhileIndexedT
//
// predicateFn is of a type "func(int,TSource)bool"
//
func (q Query) TakeWhileIndexedT(predicateFn interface{}) Query {

	whereFunc, err := newGenericFunc(
		"TakeWhileIndexedT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(index int, item interface{}) bool {
		return whereFunc.Call(index, item).(bool)
	}

	return q.TakeWhileIndexed(predicateFunc)
}
