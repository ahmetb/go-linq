package linq

// Join correlates the elements of two collections based on matching keys.
//
// A join refers to the operation of correlating the elements of two sources of
// information based on a common key. Join brings the two information sources
// and the keys by which they are matched together in one method call. This
// differs from the use of SelectMany, which requires more than one method call
// to perform the same operation.
//
// Join is a generic method: the inner element type TInner, the key type TKey, and the
// result type TResult are all inferred from the supplied functions. The key type TKey
// must be comparable.
//
// Join preserves the order of the elements of outer collection, and for each of
// these elements, the order of the matching elements of inner.
func (q Query[T]) Join[TInner any, TKey comparable, TResult any](inner Query[TInner],
	outerKeySelector func(T) TKey,
	innerKeySelector func(TInner) TKey,
	resultSelector func(outer T, inner TInner) TResult) Query[TResult] {

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
