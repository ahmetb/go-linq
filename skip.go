package linq

// Skip bypasses a specified number of elements in a collection and then returns
// the remaining elements.
func (q Query) Skip(count int) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			n := count
			q.Iterate(func(item any) bool {
				if n > 0 {
					n--
					return true
				}
				return yield(item)
			})
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
func (q Query) SkipWhile(predicate func(any) bool) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			skipping := true
			q.Iterate(func(item any) bool {
				if skipping {
					if predicate(item) {
						return true
					}
					skipping = false
				}

				return yield(item)
			})
		},
	}
}

// SkipWhileT is the typed version of SkipWhile.
//
//   - predicateFn is of type "func(TSource)bool"
//
// NOTE: SkipWhile has better performance than SkipWhileT.
func (q Query) SkipWhileT(predicateFn any) Query {

	predicateGenericFunc, err := newGenericFunc(
		"SkipWhileT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
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
func (q Query) SkipWhileIndexed(predicate func(int, any) bool) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			skipping := true
			index := 0
			q.Iterate(func(item any) bool {
				if skipping {
					if predicate(index, item) {
						index++
						return true
					}
					skipping = false
				}

				return yield(item)
			})
		},
	}
}

// SkipWhileIndexedT is the typed version of SkipWhileIndexed.
//
//   - predicateFn is of type "func(int,TSource)bool"
//
// NOTE: SkipWhileIndexed has better performance than SkipWhileIndexedT.
func (q Query) SkipWhileIndexedT(predicateFn any) Query {
	predicateGenericFunc, err := newGenericFunc(
		"SkipWhileIndexedT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(index int, item any) bool {
		return predicateGenericFunc.Call(index, item).(bool)
	}

	return q.SkipWhileIndexed(predicateFunc)
}
