package linq

// OrderByT is the typed version of OrderBy.
//
// NOTE: OrderBy method has better performance than OrderByT
//
// selectorFn is of a type "func(TSource) TKey"
//
func (q Query) OrderByT(selectorFn interface{}) OrderedQuery {

	selectorGenericFunc, err := newGenericFunc(
		"OrderByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.OrderBy(selectorFunc)

}

// OrderByDescendingT is the typed version of OrderByDescending.
//
// NOTE: OrderByDescending method has better performance than OrderByDescendingT
//
// selectorFn is of a type "func(TSource) TKey"
//
func (q Query) OrderByDescendingT(selectorFn interface{}) OrderedQuery {

	selectorGenericFunc, err := newGenericFunc(
		"OrderByDescendingT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.OrderByDescending(selectorFunc)
}

// ThenByT is the typed version of ThenBy.
//
// NOTE: ThenBy method has better performance than ThenByT
//
// selectorFn is of a type "func(TSource) TKey"
//
func (oq OrderedQuery) ThenByT(selectorFn interface{}) OrderedQuery {
	selectorGenericFunc, err := newGenericFunc(
		"ThenByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return oq.ThenBy(selectorFunc)
}

// ThenByDescendingT is the typed version of ThenByDescending.
//
// NOTE: ThenByDescending method has better performance than ThenByDescendingT
//
// selectorFn is of a type "func(TSource) TKey"
//
func (oq OrderedQuery) ThenByDescendingT(selectorFn interface{}) OrderedQuery {

	selectorFunc, ok := selectorFn.(func(interface{}) interface{})
	if !ok {
		selectorGenericFunc, err := newGenericFunc(
			"ThenByDescending", "selectorFn", selectorFn,
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
		)
		if err != nil {
			panic(err)
		}

		selectorFunc = func(item interface{}) interface{} {
			return selectorGenericFunc.Call(item)
		}
	}

	return oq.ThenByDescending(selectorFunc)
}

// SortT is the typed version of Sort.
//
// NOTE: Sort method has better performance than SortT
//
// lessFn is of a type "func(TSource,TSource) bool"
//
func (q Query) SortT(lessFn interface{}) Query {

	lessGenericFunc, err := newGenericFunc(
		"SortT", "lessFn", lessFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	lessFunc := func(i, j interface{}) bool {
		return lessGenericFunc.Call(i, j).(bool)
	}

	return q.Sort(lessFunc)
}
