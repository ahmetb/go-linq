package linq

// Reverse inverts the order of the elements in a collection.
//
// Unlike OrderBy, this sorting method does not consider the actual values
// themselves in determining the order. Rather, it just returns the elements in
// the reverse order from which they are produced by the underlying source.
func (q Query) Reverse() Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()

			items := []interface{}{}
			for item, ok := next(); ok; item, ok = next() {
				items = append(items, item)
			}

			index := len(items) - 1
			return func() (item interface{}, ok bool) {
				if index < 0 {
					return
				}

				item, ok = items[index], true
				index--
				return
			}
		},
	}
}
