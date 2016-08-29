package linq

// Select projects each element of a collection into a new form.
// Returns a query with the result of invoking the transform function
// on each element of original source.
//
// This projection method requires the transform function, selector,
// to produce one value for each value in the source collection.
// If selector returns a value that is itself a collection,
// it is up to the consumer to traverse the subcollections manually.
// In such a situation, it might be better for your query to return a single
// coalesced collection of values. To achieve this, use the SelectMany method
// instead of Select. Although SelectMany works similarly to Select,
// it differs in that the transform function returns a collection
// that is then expanded by SelectMany before it is returned.
func (q Query) Select(selector func(interface{}) interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()

			return func() (item interface{}, ok bool) {
				var it interface{}
				it, ok = next()
				if ok {
					item = selector(it)
				}

				return
			}
		},
	}
}

// SelectIndexed projects each element of a collection into a new form
// by incorporating the element's index. Returns a query with the result
// of invoking the transform function on each element of original source.
//
// The first argument to selector represents the zero-based index of that element
// in the source collection. This can be useful if the elements are in a known order
// and you want to do something with an element at a particular index,
// for example. It can also be useful if you want to retrieve the index of one
// or more elements. The second argument to selector represents the element to process.
//
// This projection method requires the transform function, selector,
// to produce one value for each value in the source collection.
// If selector returns a value that is itself a collection,
// it is up to the consumer to traverse the subcollections manually.
// In such a situation, it might be better for your query to return a single
// coalesced collection of values. To achieve this, use the SelectMany method
// instead of Select. Although SelectMany works similarly to Select,
// it differs in that the transform function returns a collection
// that is then expanded by SelectMany before it is returned.
func (q Query) SelectIndexed(selector func(int, interface{}) interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			index := 0

			return func() (item interface{}, ok bool) {
				var it interface{}
				it, ok = next()
				if ok {
					item = selector(index, it)
					index++
				}

				return
			}
		},
	}
}
