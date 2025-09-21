package linq

// Group is a type used to store the result of GroupBy method.
type Group struct {
	Key   any
	Group []any
}

// GroupBy method groups the elements of a collection according to a specified
// key selector function and projects the elements for each group by using a
// specified function.
func (q Query) GroupBy(keySelector func(any) any,
	elementSelector func(any) any) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			groups := make(map[any][]any)

			for item := range q.Iterate {
				key := keySelector(item)
				element := elementSelector(item)
				groups[key] = append(groups[key], element)
			}

			for key, group := range groups {
				group := Group{
					Key:   key,
					Group: group,
				}
				if !yield(group) {
					return
				}
			}
		},
	}
}

// GroupByT is the typed version of GroupBy.
//
//   - keySelectorFn is of type "func(TSource) TKey"
//   - elementSelectorFn is of type "func(TSource) TElement"
//
// NOTE: GroupBy has better performance than GroupByT.
func (q Query) GroupByT(keySelectorFn any,
	elementSelectorFn any) Query {
	keySelectorGenericFunc, err := newGenericFunc(
		"GroupByT", "keySelectorFn", keySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	keySelectorFunc := func(item any) any {
		return keySelectorGenericFunc.Call(item)
	}

	elementSelectorGenericFunc, err := newGenericFunc(
		"GroupByT", "elementSelectorFn", elementSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	elementSelectorFunc := func(item any) any {
		return elementSelectorGenericFunc.Call(item)

	}

	return q.GroupBy(keySelectorFunc, elementSelectorFunc)
}
