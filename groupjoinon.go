package linq

import "reflect"

// GroupJoinOn correlates the elements of two collections based on a predicate,
// and groups the results.
//
// GroupJoinOn is a more general form of GroupJoin in which the predicate
// function can be thought of as
// outerKeySelector(outerItem) == innerKeySelector(innerItem)
//
// This method produces hierarchical results, which means that elements from
// outer query are paired with collections of matching elements from inner.
// GroupJoinOn enables you to base your results on a whole set of matches for each
// element of outer query.
//
// The resultSelector function is called only one time for each outer element
// together with a collection of all the inner elements that match the outer
// element. This differs from the JoinOn method, in which the result selector
// function is invoked on pairs that contain one element from outer and one
// element from inner.
//
// GroupJoinOn preserves the order of the elements of outer, and for each element
// of outer, the order of the matching elements from inner.
func (q Query) GroupJoinOn(inner Query,
	onPredicate func(outer interface{}, inner interface{}) bool,
	resultSelector func(outer interface{}, inners []interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()

			return func() (item interface{}, ok bool) {
				outerItem, outerOk := outernext()
				if !outerOk {
					return
				}
				innernext := inner.Iterate()
				var group []interface{}
				for innerItem, ok := innernext(); ok; innerItem, ok = innernext() {
					if onPredicate(outerItem, innerItem) {
						group = append(group, innerItem)
					}
				}
				item = resultSelector(outerItem, group)
				return item, true
			}
		},
	}
}

// GroupJoinOnT is the typed version of GroupJoinOn.
//
//   - inner: The query to join to the outer query.
//   - onPredicateFn is of type "func(TOuter, TInner) bool"
//   - resultSelectorFn: is of type "func(TOuter, inners []TInner) TResult"
//
// NOTE: GroupJoinOn has better performance than GroupJoinOnT.
func (q Query) GroupJoinOnT(inner Query,
	onPredicateFn interface{},
	resultSelectorFn interface{}) Query {
	onPredicateGenericFunc, err := newGenericFunc(
		"GroupJoinOnT", "onPredicateFn", onPredicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	onPredicateFunc := func(outerItem interface{}, innerItem interface{}) bool {
		return onPredicateGenericFunc.Call(outerItem, innerItem).(bool)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"GroupJoinOnT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer interface{}, inners []interface{}) interface{} {
		innerSliceType := reflect.MakeSlice(resultSelectorGenericFunc.Cache.TypesIn[1], 0, 0)
		innersSlicePointer := reflect.New(innerSliceType.Type())
		From(inners).ToSlice(innersSlicePointer.Interface())
		innersTyped := reflect.Indirect(innersSlicePointer).Interface()
		return resultSelectorGenericFunc.Call(outer, innersTyped)
	}

	return q.GroupJoinOn(inner, onPredicateFunc, resultSelectorFunc)
}
