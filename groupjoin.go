package linq

import "reflect"

// GroupJoin correlates the elements of two collections based on key equality
// and groups the results.
//
// This method produces hierarchical results, which means that elements from
// an outer query are paired with collections of matching elements from the inner.
// GroupJoin enables you to base your results on a whole set of matches for each
// element of the outer query.
//
// The resultSelector function is called only one time for each outer element
// together with a collection of all the inner elements that match the outer
// element. This differs from the Join method, in which the result selector
// function is invoked on pairs that contain one element from outer and one
// element from inner.
//
// GroupJoin preserves the order of the elements of outer, and for each element
// of outer, the order of the matching elements from inner.
func (q Query) GroupJoin(inner Query,
	outerKeySelector func(any) any,
	innerKeySelector func(any) any,
	resultSelector func(outer any, inners []any) any) Query {

	return Query{
		Iterate: func(yield func(any) bool) {
			innerLookup := make(map[any][]any)
			for innerItem := range inner.Iterate {
				innerKey := innerKeySelector(innerItem)
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
			}

			q.Iterate(func(outerItem any) bool {
				outerKey := outerKeySelector(outerItem)
				innerGroup, ok := innerLookup[outerKey]

				var result any
				if ok {
					result = resultSelector(outerItem, innerGroup)
				} else {
					result = resultSelector(outerItem, []any{})
				}

				return yield(result)
			})
		},
	}
}

// GroupJoinT is the typed version of GroupJoin.
//
//   - inner: The query to join to the outer query.
//   - outerKeySelectorFn is of type "func(TOuter) TKey"
//   - innerKeySelectorFn is of type "func(TInner) TKey"
//   - resultSelectorFn: is of type "func(TOuter, inners []TInner) TResult"
//
// NOTE: GroupJoin has better performance than GroupJoinT.
func (q Query) GroupJoinT(inner Query,
	outerKeySelectorFn any,
	innerKeySelectorFn any,
	resultSelectorFn any) Query {
	outerKeySelectorGenericFunc, err := newGenericFunc(
		"GroupJoinT", "outerKeySelectorFn", outerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	outerKeySelectorFunc := func(item any) any {
		return outerKeySelectorGenericFunc.Call(item)
	}

	innerKeySelectorFuncGenericFunc, err := newGenericFunc(
		"GroupJoinT", "innerKeySelectorFn", innerKeySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	innerKeySelectorFunc := func(item any) any {
		return innerKeySelectorFuncGenericFunc.Call(item)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"GroupJoinT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer any, inners []any) any {
		innerSliceType := reflect.MakeSlice(resultSelectorGenericFunc.Cache.TypesIn[1], 0, 0)
		innersSlicePointer := reflect.New(innerSliceType.Type())
		From(inners).ToSlice(innersSlicePointer.Interface())
		innersTyped := reflect.Indirect(innersSlicePointer).Interface()
		return resultSelectorGenericFunc.Call(outer, innersTyped)
	}

	return q.GroupJoin(inner, outerKeySelectorFunc, innerKeySelectorFunc, resultSelectorFunc)
}
