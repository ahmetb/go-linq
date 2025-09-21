package linq

// Intersect produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
func (q Query) Intersect(q2 Query) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			set := make(map[any]struct{})
			for item := range q2.Iterate {
				set[item] = struct{}{}
			}

			for item := range q.Iterate {
				if _, exists := set[item]; exists {
					delete(set, item)
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

// IntersectBy produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
//
// IntersectBy invokes a transform function on each element of both collections.
func (q Query) IntersectBy(q2 Query,
	selector func(any) any) Query {

	return Query{
		Iterate: func(yield func(any) bool) {
			set := make(map[any]struct{})
			for item := range q2.Iterate {
				key := selector(item)
				set[key] = struct{}{}
			}

			for item := range q.Iterate {
				key := selector(item)
				if _, exists := set[key]; exists {
					delete(set, key)
					if !yield(item) {
						return
					}
				}
			}
		},
	}
}

// IntersectByT is the typed version of IntersectBy.
//
//   - selectorFn is of type "func(TSource) TSource"
//
// NOTE: IntersectBy has better performance than IntersectByT.
func (q Query) IntersectByT(q2 Query,
	selectorFn any) Query {
	selectorGenericFunc, err := newGenericFunc(
		"IntersectByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item any) any {
		return selectorGenericFunc.Call(item)
	}

	return q.IntersectBy(q2, selectorFunc)
}
