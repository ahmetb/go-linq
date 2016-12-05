package linq

// SkipWhileT is the typed version of SkipWhile.
//
// NOTE: SkipWhile method has better performance than SkipWhileT
//
// predicateFn is of a type "func(TSource)bool"
//
func (q Query) SkipWhileT(predicateFn interface{}) Query {

	predicateGenericFunc, err := newGenericFunc(
		"SkipWhileT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.SkipWhile(predicateFunc)
}

// SkipWhileIndexedT is the typed version of SkipWhileIndexed.
//
// NOTE: SkipWhileIndexed method has better performance than SkipWhileIndexedT
//
// predicateFn is of a type "func(int,TSource)bool"
//
func (q Query) SkipWhileIndexedT(predicateFn interface{}) Query {

	predicateGenericFunc, err := newGenericFunc(
		"SkipWhileIndexedT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(index int, item interface{}) bool {
		return predicateGenericFunc.Call(index, item).(bool)
	}

	return q.SkipWhileIndexed(predicateFunc)
}
