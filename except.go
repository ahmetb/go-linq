package linq

// Except produces the set difference of two sequences. The set difference is
// the members of the first sequence that don't appear in the second sequence.
func (q Query) Except(q2 Query) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			set := make(map[any]struct{})
			for item := range q2.Iterate {
				set[item] = struct{}{}
			}

			q.Iterate(func(item any) bool {
				if _, seen := set[item]; !seen {
					return yield(item)
				}
				return true
			})
		},
	}
}

// ExceptBy invokes a transform function on each element of a collection and
// produces the set difference of two sequences. The set difference is the
// members of the first sequence that don't appear in the second sequence.
func (q Query) ExceptBy(q2 Query, selector func(any) any) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			set := make(map[any]struct{})
			for item := range q2.Iterate {
				key := selector(item)
				set[key] = struct{}{}
			}

			q.Iterate(func(item any) bool {
				key := selector(item)
				if _, seen := set[key]; !seen {
					return yield(item)
				}
				return true
			})
		},
	}
}

// ExceptByT is the typed version of ExceptBy.
//
//   - selectorFn is of type "func(TSource) TSource"
//
// NOTE: ExceptBy has better performance than ExceptByT.
func (q Query) ExceptByT(q2 Query,
	selectorFn any) Query {
	selectorGenericFunc, err := newGenericFunc(
		"ExceptByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item any) any {
		return selectorGenericFunc.Call(item)
	}

	return q.ExceptBy(q2, selectorFunc)
}
