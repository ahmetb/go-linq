package linq

// Where filters a collection of values based on a predicate.
func (q Query) Where(predicate func(any) bool) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			q.Iterate(func(item any) bool {
				if predicate(item) {
					return yield(item)
				}
				return true
			})
		},
	}
}

// WhereT is the typed version of Where.
//
//   - predicateFn is of type "func(TSource)bool"
//
// NOTE: Where has better performance than WhereT.
func (q Query) WhereT(predicateFn any) Query {

	predicateGenericFunc, err := newGenericFunc(
		"WhereT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.Where(predicateFunc)
}

// WhereIndexed filters a collection of values based on a predicate. Each
// element's index is used in the logic of the predicate function.
//
// The first argument represents the zero-based index of the element within
// the collection. The second argument of predicate represents the element to test.
func (q Query) WhereIndexed(predicate func(int, any) bool) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			index := 0
			q.Iterate(func(item any) bool {
				shouldYield := predicate(index, item)
				index++

				if shouldYield {
					return yield(item)
				}

				return true
			})
		},
	}
}

// WhereIndexedT is the typed version of WhereIndexed.
//
//   - predicateFn is of type "func(int,TSource)bool"
//
// NOTE: WhereIndexed has better performance than WhereIndexedT.
func (q Query) WhereIndexedT(predicateFn any) Query {
	predicateGenericFunc, err := newGenericFunc(
		"WhereIndexedT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(index int, item any) bool {
		return predicateGenericFunc.Call(index, item).(bool)
	}

	return q.WhereIndexed(predicateFunc)
}
