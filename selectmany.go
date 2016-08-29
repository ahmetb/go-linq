package linq

// SelectMany projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection.
func (q Query) SelectMany(selector func(interface{}) Query) Query {
	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			var inner interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if inner == nil {
						inner, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(inner).Iterate()
					}

					item, ok = innernext()
					if !ok {
						inner = nil
					}
				}

				return
			}
		},
	}
}

// SelectManyIndexed projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection.
//
// The first argument to selector represents the zero-based index of that element
// in the source collection. This can be useful if the elements are in a known order
// and you want to do something with an element at a particular index, for example.
// It can also be useful if you want to retrieve the index of one or more elements.
// The second argument to selector represents the element to process.
func (q Query) SelectManyIndexed(selector func(int, interface{}) Query) Query {
	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			index := 0
			var inner interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if inner == nil {
						inner, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(index, inner).Iterate()
						index++
					}

					item, ok = innernext()
					if !ok {
						inner = nil
					}
				}

				return
			}
		},
	}
}

// SelectManyBy projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection, and invokes
// a result selector function on each element therein.
func (q Query) SelectManyBy(
	selector func(interface{}) Query,
	resultSelector func(interface{}, interface{}) interface{},
) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			var outer interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if outer == nil {
						outer, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(outer).Iterate()
					}

					item, ok = innernext()
					if !ok {
						outer = nil
					}
				}

				item = resultSelector(outer, item)
				return
			}
		},
	}
}

// SelectManyByIndexed projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection, and invokes
// a result selector function on each element therein.
// The index of each source element is used in the intermediate projected form
// of that element.
func (q Query) SelectManyByIndexed(selector func(int, interface{}) Query,
	resultSelector func(interface{}, interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			index := 0
			var outer interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if outer == nil {
						outer, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(index, outer).Iterate()
						index++
					}

					item, ok = innernext()
					if !ok {
						outer = nil
					}
				}

				item = resultSelector(outer, item)
				return
			}
		},
	}
}
