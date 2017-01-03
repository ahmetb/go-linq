package linq

// Distinct method returns distinct elements from a collection. The result is an
// unordered collection that contains no duplicate values.
func (q Query) Distinct() Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			set := make(map[interface{}]bool)

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if _, has := set[item]; !has {
						set[item] = true
						return
					}
				}

				return
			}
		},
	}
}

// Distinct method returns distinct elements from a collection. The result is an
// ordered collection that contains no duplicate values.
//
// NOTE: Distinct method on OrderedQuery type has better performance than
// Distinct method on Query type.
func (oq OrderedQuery) Distinct() OrderedQuery {
	return OrderedQuery{
		orders: oq.orders,
		Query: Query{
			Iterate: func() Iterator {
				next := oq.Iterate()
				var prev interface{}

				return func() (item interface{}, ok bool) {
					for item, ok = next(); ok; item, ok = next() {
						if item != prev {
							prev = item
							return
						}
					}

					return
				}
			},
		},
	}
}

// DistinctBy method returns distinct elements from a collection. This method
// executes selector function for each element to determine a value to compare.
// The result is an unordered collection that contains no duplicate values.
func (q Query) DistinctBy(selector func(interface{}) interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			set := make(map[interface{}]bool)

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					s := selector(item)
					if _, has := set[s]; !has {
						set[s] = true
						return
					}
				}

				return
			}
		},
	}
}

// DistinctByT is the typed version of DistinctBy.
//
//   - selectorFn is of type "func(TSource) TSource".
//
// NOTE: DistinctBy has better performance than DistinctByT.
func (q Query) DistinctByT(selectorFn interface{}) Query {
	selectorFunc, ok := selectorFn.(func(interface{}) interface{})
	if !ok {
		selectorGenericFunc, err := newGenericFunc(
			"DistinctByT", "selectorFn", selectorFn,
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
		)
		if err != nil {
			panic(err)
		}

		selectorFunc = func(item interface{}) interface{} {
			return selectorGenericFunc.Call(item)
		}
	}
	return q.DistinctBy(selectorFunc)
}
