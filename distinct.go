package linq

// Distinct method returns distinct elements from a collection. The result is an
// unordered collection that contains no duplicate values.
func (q Query) Distinct() Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			set := make(map[any]struct{})

			q.Iterate(func(item any) bool {
				if _, seen := set[item]; !seen {
					set[item] = struct{}{}
					return yield(item)
				}

				return true
			})
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
			Iterate: func(yield func(any) bool) {
				var previous any
				isFirst := true

				oq.Iterate(func(item any) bool {
					if isFirst || item != previous {
						previous = item
						isFirst = false
						return yield(item)
					}

					return true
				})
			},
		},
	}
}

// DistinctBy method returns distinct elements from a collection. This method
// executes selector function for each element to determine a value to compare.
// The result is an unordered collection that contains no duplicate values.
func (q Query) DistinctBy(selector func(any) any) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			set := make(map[any]struct{})

			q.Iterate(func(item any) bool {
				key := selector(item)

				if _, seen := set[key]; !seen {
					set[key] = struct{}{}
					return yield(item)
				}

				return true
			})
		},
	}
}

// DistinctByT is the typed version of DistinctBy.
//
//   - selectorFn is of type "func(TSource) TSource".
//
// NOTE: DistinctBy has better performance than DistinctByT.
func (q Query) DistinctByT(selectorFn any) Query {
	selectorFunc, ok := selectorFn.(func(any) any)
	if !ok {
		selectorGenericFunc, err := newGenericFunc(
			"DistinctByT", "selectorFn", selectorFn,
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
		)
		if err != nil {
			panic(err)
		}

		selectorFunc = func(item any) any {
			return selectorGenericFunc.Call(item)
		}
	}
	return q.DistinctBy(selectorFunc)
}
