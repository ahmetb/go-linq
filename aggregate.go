package linq

// Aggregate applies an accumulator function over a sequence.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The first element of source is used as the initial aggregate value. The
// result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (q Query) Aggregate(f func(interface{}, interface{}) interface{}) interface{} {
	next := q.Iterate()

	result, any := next()
	if !any {
		return nil
	}

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return result
}

// AggregateT is the typed version of Aggregate.
//
//   - f is of type: func(TSource, TSource) TSource
//
// NOTE: Aggregate has better performance than AggregateT.
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

// AggregateWithSeed applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The value of the seed parameter is used as the initial aggregate value.
// The result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (q Query) AggregateWithSeed(seed interface{},
	f func(interface{}, interface{}) interface{}) interface{} {

	next := q.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return result
}

// AggregateWithSeedT is the typed version of AggregateWithSeed.
//
//   - f is of type "func(TAccumulate, TSource) TAccumulate"
//
// NOTE: AggregateWithSeed has better performance than
// AggregateWithSeedT.
func (q Query) AggregateWithSeedT(seed interface{},
	f interface{}) interface{} {
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

// AggregateWithSeedBy applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value, and the
// specified function is used to select the result value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source.
// Each time func is called, Aggregate passes both the element from the sequence
// and an aggregated value (as the first argument to func). The value of the
// seed parameter is used as the initial aggregate value. The result of func
// replaces the previous aggregated value.
//
// The final result of func is passed to resultSelector to obtain the final
// result of Aggregate.
func (q Query) AggregateWithSeedBy(seed interface{},
	f func(interface{}, interface{}) interface{},
	resultSelector func(interface{}) interface{}) interface{} {

	next := q.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return resultSelector(result)
}

// AggregateWithSeedByT is the typed version of AggregateWithSeedBy.
//
//   - f is of type "func(TAccumulate, TSource) TAccumulate"
//   - resultSelectorFn is of type "func(TAccumulate) TResult"
//
// NOTE: AggregateWithSeedBy has better performance than
// AggregateWithSeedByT.
func (q Query) AggregateWithSeedByT(seed interface{},
	f interface{},
	resultSelectorFn interface{}) interface{} {
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
