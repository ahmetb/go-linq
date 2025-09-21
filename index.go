package linq

// IndexOf searches for an element that matches the conditions defined by a specified predicate
// and returns the zero-based index of the first occurrence within the collection. This method
// returns -1 if an item that matches the conditions is not found.
func (q Query) IndexOf(predicate func(any) bool) int {
	index := 0
	for item := range q.Iterate {
		if predicate(item) {
			return index
		}
		index++
	}

	return -1
}

// IndexOfT is the typed version of IndexOf.
//
//   - predicateFn is of type "func(int,TSource)bool"
//
// NOTE: IndexOf has better performance than IndexOfT.
func (q Query) IndexOfT(predicateFn any) int {

	predicateGenericFunc, err := newGenericFunc(
		"IndexOfT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.IndexOf(predicateFunc)
}
