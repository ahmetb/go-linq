package linq

// Join correlates the elements of two collection based on matching keys.
//
// A join refers to the operation of correlating the elements of two sources of
// information based on a common key. Join brings the two information sources
// and the keys by which they are matched together in one method call. This
// differs from the use of SelectMany, which requires more than one method call
// to perform the same operation.
//
// Join preserves the order of the elements of outer collection, and for each of
// these elements, the order of the matching elements of inner.
func (q Query) Join(inner Query,
	outerKeySelector func(interface{}) interface{},
	innerKeySelector func(interface{}) interface{},
	resultSelector func(outer interface{}, inner interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			innernext := inner.Iterate()

			innerLookup := make(map[interface{}][]interface{})
			for innerItem, ok := innernext(); ok; innerItem, ok = innernext() {
				innerKey := innerKeySelector(innerItem)
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
			}

			var outerItem interface{}
			var innerGroup []interface{}
			innerLen, innerIndex := 0, 0

			return func() (item interface{}, ok bool) {
				if innerIndex >= innerLen {
					has := false
					for !has {
						outerItem, ok = outernext()
						if !ok {
							return
						}

						innerGroup, has = innerLookup[outerKeySelector(outerItem)]
						innerLen = len(innerGroup)
						innerIndex = 0
					}
				}

				item = resultSelector(outerItem, innerGroup[innerIndex])
				innerIndex++
				return item, true
			}
		},
	}
}

// JoinT is the typed version of Join.
//
//   - outerKeySelectorFn is of type "func(TOuter) TKey"
//   - innerKeySelectorFn is of type "func(TInner) TKey"
//   - resultSelectorFn is of type "func(TOuter,TInner) TResult"
//
// NOTE: Join has better performance than JoinT.
func (q Query) JoinT(inner Query,
	outerKeySelectorFn interface{},
	innerKeySelectorFn interface{},
	resultSelectorFn interface{}) Query {
	outerKeySelectorGenericFunc, err := newGenericFunc(
		"JoinT", "outerKeySelectorFn", outerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	outerKeySelectorFunc := func(item interface{}) interface{} {
		return outerKeySelectorGenericFunc.Call(item)
	}

	innerKeySelectorFuncGenericFunc, err := newGenericFunc(
		"JoinT", "innerKeySelectorFn",
		innerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	innerKeySelectorFunc := func(item interface{}) interface{} {
		return innerKeySelectorFuncGenericFunc.Call(item)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"JoinT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer interface{}, inner interface{}) interface{} {
		return resultSelectorGenericFunc.Call(outer, inner)
	}

	return q.Join(inner, outerKeySelectorFunc, innerKeySelectorFunc, resultSelectorFunc)
}
