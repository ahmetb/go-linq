package linq

// JoinOn correlates the elements of two collections based on a predicate function.
//
// A joinOn refers to the operation of correlating the elements of two sources of
// information based on a predicate function that returns a bool given a pair of
// outer and inner collection elements. JoinOn is a more general form of Join
// in which the predicate function can be thought of as
// outerKeySelector(outerItem) == innerKeySelector(innerItem)
//
// JoinOn preserves the order of the elements of outer collection, and for each of
// these elements, the order of the matching elements of inner.
func (q Query) JoinOn(inner Query,
	onPredicate func(outer interface{}, inner interface{}) bool,
	resultSelector func(outer interface{}, inner interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			innernext := inner.Iterate()

			outerItem, outerOk := outernext()

			return func() (item interface{}, ok bool) {

				for outerOk {
					for innerItem, ok := innernext(); ok; innerItem, ok = innernext() {
						if onPredicate(outerItem, innerItem) {
							item = resultSelector(outerItem, innerItem)
							return item, true
						}
					}
					innernext = inner.Iterate()
					outerItem, outerOk = outernext()
				}
				return
			}
		},
	}
}

// JoinOnT is the typed version of JoinOn.
//
//   - inner: The query to join to the outer query.
//   - onPredicateFn is of type "func(TOuter, TInner) bool"
//   - resultSelectorFn is of type "func(TOuter,TInner) TResult"
//
// NOTE: JoinOn has better performance than JoinOnT.
func (q Query) JoinOnT(inner Query,
	onPredicateFn interface{},
	resultSelectorFn interface{}) Query {
	onPredicateGenericFunc, err := newGenericFunc(
		"JoinOnT", "onPredicateFn", onPredicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	onPredicateFunc := func(outerItem interface{}, innerItem interface{}) bool {
		return onPredicateGenericFunc.Call(outerItem, innerItem).(bool)
	}

	resultSelectorGenericFunc, err := newGenericFunc(
		"JoinOnT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(outer interface{}, inner interface{}) interface{} {
		return resultSelectorGenericFunc.Call(outer, inner)
	}

	return q.JoinOn(inner, onPredicateFunc, resultSelectorFunc)
}
