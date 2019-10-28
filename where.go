package linq

// Where filters a collection of values based on a predicate.
func (q Query) Where(predicate func(interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if predicate(item) {
						return
					}
				}

				return
			}
		},
	}
}

// WhereT is the typed version of Where.
//
//   - predicateFn is of type "func(TSource)bool"
//
// NOTE: Where has better performance than WhereT.
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

// WhereIndexed filters a collection of values based on a predicate. Each
// element's index is used in the logic of the predicate function.
//
// The first argument represents the zero-based index of the element within
// collection. The second argument of predicate represents the element to test.
func (q Query) WhereIndexed(predicate func(int, interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			index := 0

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if predicate(index, item) {
						index++
						return
					}

					index++
				}

				return
			}
		},
	}
}

// WhereIndexedT is the typed version of WhereIndexed.
//
//   - predicateFn is of type "func(int,TSource)bool"
//
// NOTE: WhereIndexed has better performance than WhereIndexedT.
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
