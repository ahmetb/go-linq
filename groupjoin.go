package linq

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
// GroupJoin is a generic method: the inner element type TInner, the key type TKey, and
// the result type TResult are all inferred from the supplied functions. The key type
// TKey must be comparable.
//
// GroupJoin preserves the order of the elements of outer, and for each element
// of outer, the order of the matching elements from inner.
func (q Query[T]) GroupJoin[TInner any, TKey comparable, TResult any](inner Query[TInner],
	outerKeySelector func(T) TKey,
	innerKeySelector func(TInner) TKey,
	resultSelector func(outer T, inners []TInner) TResult) Query[TResult] {

	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			innerLookup := make(map[TKey][]TInner)
			inner.Iterate(func(innerItem TInner) bool {
				innerKey := innerKeySelector(innerItem)
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
				return true
			})

			q.Iterate(func(outerItem T) bool {
				outerKey := outerKeySelector(outerItem)
				innerGroup, ok := innerLookup[outerKey]
				if !ok {
					innerGroup = []TInner{}
				}

				return yield(resultSelector(outerItem, innerGroup))
			})
		},
	}
}
