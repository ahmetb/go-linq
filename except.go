package linq

// Except produces the set difference of two sequences. The set difference is
// the members of the first sequence that don't appear in the second sequence.
func (q Query) Except(q2 Query) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()

			next2 := q2.Iterate()
			set := make(map[interface{}]bool)
			for i, ok := next2(); ok; i, ok = next2() {
				set[i] = true
			}

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if _, has := set[item]; !has {
						return
					}
				}

				return
			}
		},
	}
}

// ExceptBy invokes a transform function on each element of a collection and
// produces the set difference of two sequences. The set difference is the
// members of the first sequence that don't appear in the second sequence.
func (q Query) ExceptBy(q2 Query,
	selector func(interface{}) interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()

			next2 := q2.Iterate()
			set := make(map[interface{}]bool)
			for i, ok := next2(); ok; i, ok = next2() {
				s := selector(i)
				set[s] = true
			}

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					s := selector(item)
					if _, has := set[s]; !has {
						return
					}
				}

				return
			}
		},
	}
}

// ExceptByT is the typed version of ExceptBy.
//
//   - selectorFn is of type "func(TSource) TSource"
//
// NOTE: ExceptBy has better performance than ExceptByT.
func (q Query) ExceptByT(q2 Query,
	selectorFn interface{}) Query {
	selectorGenericFunc, err := newGenericFunc(
		"ExceptByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.ExceptBy(q2, selectorFunc)
}
