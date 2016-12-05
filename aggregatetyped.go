package linq

// AggregateT is the typed version of Aggregate.
//
// NOTE: Aggregate method has better performance than AggregateT
//
// f is of type: func(TSource, TSource) TSource
//
func (q Query) AggregateT(f interface{}) interface{} {

	fGenericFunc, err := newGenericFunc(
		"AggregateT", "f", f,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	fFunc := func(result interface{}, current interface{}) interface{} {
		return fGenericFunc.Call(result, current)
	}

	return q.Aggregate(fFunc)
}

// AggregateWithSeedT is the typed version of AggregateWithSeed.
//
// NOTE: AggregateWithSeed method has better performance than AggregateWithSeedT
//
// f is of a type "func(TAccumulate, TSource) TAccumulate"
//
func (q Query) AggregateWithSeedT(seed interface{}, f interface{}) interface{} {

	fGenericFunc, err := newGenericFunc(
		"AggregateWithSeed", "f", f,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	fFunc := func(result interface{}, current interface{}) interface{} {
		return fGenericFunc.Call(result, current)
	}

	return q.AggregateWithSeed(seed, fFunc)
}

// AggregateWithSeedByT is the typed version of AggregateWithSeedBy.
//
// NOTE: AggregateWithSeedBy method has better performance than AggregateWithSeedByT
//
// f is of a type "func(TAccumulate, TSource) TAccumulate"
//
// resultSelectorFn is of type "func(TAccumulate) TResult"
//
func (q Query) AggregateWithSeedByT(
	seed interface{},
	f interface{},
	resultSelectorFn interface{},
) interface{} {

	fGenericFunc, err := newGenericFunc(
		"AggregateWithSeedByT", "f", f,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	fFunc := func(result interface{}, current interface{}) interface{} {
		return fGenericFunc.Call(result, current)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"AggregateWithSeedByT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(result interface{}) interface{} {
		return resultSelectorGenericFunc.Call(result)
	}

	return q.AggregateWithSeedBy(seed, fFunc, resultSelectorFunc)
}
