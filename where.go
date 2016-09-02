package linq

// Where filters a collection of values based on a predicate.
func (q Query) Where(predicate func(interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if predicate(item) {
						return
					}
				}

				return
			}
		},
	}
}

// WhereIndexed filters a collection of values based on a predicate.
// Each element's index is used in the logic of the predicate function.
//
// The first argument represents the zero-based index of the element within collection.
// The second argument of predicate represents the element to test.
func (q Query) WhereIndexed(predicate func(int, interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			index := 0

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if predicate(index, item) {
						return
					}

					index++
				}

				return
			}
		},
	}
}
