package linq

import "iter"

// Zip applies a specified function to the corresponding elements of two
// collections, producing a collection of the results.
//
// The method steps through the two input collections, applying function
// resultSelector to corresponding elements of the two collections. The method
// returns a collection of the values that are returned by resultSelector. If
// the input collections do not have the same number of elements, the method
// combines elements until it reaches the end of one of the collections. For
// example, if one collection has three elements and the other one has four, the
// result collection has only three elements.
func (q Query) Zip(q2 Query,
	resultSelector func(any, any) any) Query {

	return Query{
		Iterate: func(yield func(any) bool) {
			next1, stop1 := iter.Pull(q.Iterate)
			defer stop1()

			next2, stop2 := iter.Pull(q2.Iterate)
			defer stop2()

			for {
				item1, ok1 := next1()
				item2, ok2 := next2()

				if !ok1 || !ok2 {
					return
				}

				result := resultSelector(item1, item2)
				if !yield(result) {
					return
				}
			}
		},
	}
}

// ZipT is the typed version of Zip.
//
//   - resultSelectorFn is of type "func(TFirst,TSecond)TResult"
//
// NOTE: Zip has better performance than ZipT.
func (q Query) ZipT(q2 Query,
	resultSelectorFn any) Query {
	resultSelectorGenericFunc, err := newGenericFunc(
		"ZipT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(item1 any, item2 any) any {
		return resultSelectorGenericFunc.Call(item1, item2)
	}

	return q.Zip(q2, resultSelectorFunc)
}
