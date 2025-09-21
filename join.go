package linq

// Join correlates the elements of two collections based on matching keys.
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
	outerKeySelector func(any) any,
	innerKeySelector func(any) any,
	resultSelector func(outer any, inner any) any) Query {

	return Query{
		Iterate: func(yield func(any) bool) {
			innerLookup := make(map[any][]any)
			for innerItem := range inner.Iterate {
				innerKey := innerKeySelector(innerItem)
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
			}

			q.Iterate(func(outerItem any) bool {
				outerKey := outerKeySelector(outerItem)

				if innerGroup, ok := innerLookup[outerKey]; ok {
					for _, innerItem := range innerGroup {
						result := resultSelector(outerItem, innerItem)
						if !yield(result) {
							return false
						}
					}
				}
				return true
			})
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
	outerKeySelectorFn any,
	innerKeySelectorFn any,
	resultSelectorFn any) Query {
	outerKeySelectorGenericFunc, err := newGenericFunc(
		"JoinT", "outerKeySelectorFn", outerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	outerKeySelectorFunc := func(item any) any {
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

	innerKeySelectorFunc := func(item any) any {
		return innerKeySelectorFuncGenericFunc.Call(item)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"JoinT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer any, inner any) any {
		return resultSelectorGenericFunc.Call(outer, inner)
	}

	return q.Join(inner, outerKeySelectorFunc, innerKeySelectorFunc, resultSelectorFunc)
}
