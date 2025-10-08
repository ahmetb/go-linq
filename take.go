package linq

// Take returns a specified number of contiguous elements from the start of a
// collection.
func (q Query) Take(count int) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			n := count
			q.Iterate(func(item any) bool {
				if n > 0 {
					n--
					return yield(item)
				}
				return false
			})
		},
	}
}

// TakeWhile returns elements from a collection as long as a specified condition
// is true and then skips the remaining elements.
func (q Query) TakeWhile(predicate func(any) bool) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			q.Iterate(func(item any) bool {
				if predicate(item) {
					return yield(item)
				}
				return false
			})
		},
	}
}

// TakeWhileT is the typed version of TakeWhile.
//
//   - predicateFn is of type "func(TSource)bool"
//
// NOTE: TakeWhile has better performance than TakeWhileT.
func (q Query) TakeWhileT(predicateFn any) Query {

	predicateGenericFunc, err := newGenericFunc(
		"TakeWhileT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.TakeWhile(predicateFunc)
}

// TakeWhileIndexed returns elements from a collection as long as a specified
// condition is true. The element's index is used in the logic of the predicate
// function. The first argument of predicate represents the zero-based index of
// the element within the collection. The second argument represents the element to
// test.
func (q Query) TakeWhileIndexed(predicate func(int, any) bool) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			index := 0
			q.Iterate(func(item any) bool {
				if predicate(index, item) {
					index++
					return yield(item)
				}
				return false
			})
		},
	}
}

// TakeWhileIndexedT is the typed version of TakeWhileIndexed.
//
//   - predicateFn is of type "func(int,TSource)bool"
//
// NOTE: TakeWhileIndexed has better performance than TakeWhileIndexedT.
func (q Query) TakeWhileIndexedT(predicateFn any) Query {
	whereFunc, err := newGenericFunc(
		"TakeWhileIndexedT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(index int, item any) bool {
		return whereFunc.Call(index, item).(bool)
	}

	return q.TakeWhileIndexed(predicateFunc)
}
