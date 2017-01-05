package linq

// Skip bypasses a specified number of elements in a collection and then returns
// the remaining elements.
func (q Query) Skip(count int) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			n := count

			return func() (item interface{}, ok bool) {
				for ; n > 0; n-- {
					item, ok = next()
					if !ok {
						return
					}
				}

				return next()
			}
		},
	}
}

// SkipWhile bypasses elements in a collection as long as a specified condition
// is true and then returns the remaining elements.
//
// This method tests each element by using predicate and skips the element if
// the result is true. After the predicate function returns false for an
// element, that element and the remaining elements in source are returned and
// there are no more invocations of predicate.
func (q Query) SkipWhile(predicate func(interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			ready := false

			return func() (item interface{}, ok bool) {
				for !ready {
					item, ok = next()
					if !ok {
						return
					}

					ready = !predicate(item)
					if ready {
						return
					}
				}

				return next()
			}
		},
	}
}

// SkipWhileT is the typed version of SkipWhile.
//
//   - predicateFn is of type "func(TSource)bool"
//
// NOTE: SkipWhile has better performance than SkipWhileT.
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

// SkipWhileIndexed bypasses elements in a collection as long as a specified
// condition is true and then returns the remaining elements. The element's
// index is used in the logic of the predicate function.
//
// This method tests each element by using predicate and skips the element if
// the result is true. After the predicate function returns false for an
// element, that element and the remaining elements in source are returned and
// there are no more invocations of predicate.
func (q Query) SkipWhileIndexed(predicate func(int, interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			ready := false
			index := 0

			return func() (item interface{}, ok bool) {
				for !ready {
					item, ok = next()
					if !ok {
						return
					}

					ready = !predicate(index, item)
					if ready {
						return
					}

					index++
				}

				return next()
			}
		},
	}
}

// SkipWhileIndexedT is the typed version of SkipWhileIndexed.
//
//   - predicateFn is of type "func(int,TSource)bool"
//
// NOTE: SkipWhileIndexed has better performance than SkipWhileIndexedT.
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
