package linq

// Group is a type that is used to store the result of GroupBy method.
type Group struct {
	Key   interface{}
	Group []interface{}
}

type GroupG[T1, T2 any] struct {
	Key   T1
	Group []T2
}

// GroupBy method groups the elements of a collection according to a specified
// key selector function and projects the elements for each group by using a
// specified function.
func (q Query) GroupBy(keySelector func(interface{}) interface{},
	elementSelector func(interface{}) interface{}) Query {
	return Query{
		func() Iterator {
			next := q.Iterate()
			set := make(map[interface{}][]interface{})

			for item, ok := next(); ok; item, ok = next() {
				key := keySelector(item)
				set[key] = append(set[key], elementSelector(item))
			}

			len := len(set)
			idx := 0
			groups := make([]Group, len)
			for k, v := range set {
				groups[idx] = Group{k, v}
				idx++
			}

			index := 0

			return func() (item interface{}, ok bool) {
				ok = index < len
				if ok {
					item = groups[index]
					index++
				}

				return
			}
		},
	}
}

// GroupBy method groups the elements of a collection according to a specified
// key selector function and projects the elements for each group by using a
// specified function.
func (e *Expended3[T, TK, TE]) GroupBy(keySelector func(T) TK,
	elementSelector func(T) TE) QueryG[GroupG[TK, TE]] {
	return QueryG[GroupG[TK, TE]]{
		func() IteratorG[GroupG[TK, TE]] {
			next := e.q.Iterate()
			set := make(map[interface{}][]TE)

			for item, ok := next(); ok; item, ok = next() {
				key := keySelector(item)
				set[key] = append(set[key], elementSelector(item))
			}

			length := len(set)
			idx := 0
			groups := make([]GroupG[TK, TE], length)
			for k, v := range set {
				groups[idx] = GroupG[TK, TE]{Key: k.(TK), Group: v}
				idx++
			}

			index := 0

			return func() (item GroupG[TK, TE], ok bool) {
				ok = index < length
				if ok {
					item = groups[index]
					index++
				}

				return
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
func (q Query) GroupByT(keySelectorFn interface{},
	elementSelectorFn interface{}) Query {
	keySelectorGenericFunc, err := newGenericFunc(
		"GroupByT", "keySelectorFn", keySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	keySelectorFunc := func(item interface{}) interface{} {
		return keySelectorGenericFunc.Call(item)
	}

	elementSelectorGenericFunc, err := newGenericFunc(
		"GroupByT", "elementSelectorFn", elementSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	elementSelectorFunc := func(item interface{}) interface{} {
		return elementSelectorGenericFunc.Call(item)

	}

	return q.GroupBy(keySelectorFunc, elementSelectorFunc)
}
