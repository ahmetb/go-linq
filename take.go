package linq

// Take returns a specified number of contiguous elements from the start of a
// collection.
func (q Query) Take(count int) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			n := count

			return func() (item interface{}, ok bool) {
				if n <= 0 {
					return
				}

				n--
				return next()
			}
		},
	}
}

// TakeWhile returns elements from a collection as long as a specified condition
// is true, and then skips the remaining elements.
func (q Query) TakeWhile(predicate func(interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			done := false

			return func() (item interface{}, ok bool) {
				if done {
					return
				}

				item, ok = next()
				if !ok {
					done = true
					return
				}

				if predicate(item) {
					return
				}

				done = true
				return nil, false
			}
		},
	}
}

// TakeWhileT is the typed version of TakeWhile.
//
//   - predicateFn is of type "func(TSource)bool"
//
// NOTE: TakeWhile has better performance than TakeWhileT.
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

// TakeWhileIndexed returns elements from a collection as long as a specified
// condition is true. The element's index is used in the logic of the predicate
// function. The first argument of predicate represents the zero-based index of
// the element within collection. The second argument represents the element to
// test.
func (q Query) TakeWhileIndexed(predicate func(int, interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			done := false
			index := 0

			return func() (item interface{}, ok bool) {
				if done {
					return
				}

				item, ok = next()
				if !ok {
					done = true
					return
				}

				if predicate(index, item) {
					index++
					return
				}

				done = true
				return nil, false
			}
		},
	}
}

// TakeWhileIndexedT is the typed version of TakeWhileIndexed.
//
//   - predicateFn is of type "func(int,TSource)bool"
//
// NOTE: TakeWhileIndexed has better performance than TakeWhileIndexedT.
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
