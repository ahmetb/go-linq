package linq

// Select projects each element of a collection into a new form. Returns a query
// with the result of invoking the transform function on each element of
// the original source.
//
// Select is a generic method: the result element type TResult is inferred from the
// selector function, so type-changing projections chain naturally.
//
// This projection method requires the transform function, selector, to produce
// one value for each value in the source collection. If selector returns a
// value that is itself a collection, it is up to the consumer to traverse the
// subcollections manually. In such a situation, it might be better for your
// query to return a single coalesced collection of values. To achieve this, use
// the SelectMany method instead of Select. Although SelectMany works similarly
// to Select, it differs in that the transform function returns a collection
// that is then expanded by SelectMany before it is returned.
func (q Query[T]) Select[TResult any](selector func(T) TResult) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			q.Iterate(func(item T) bool {
				return yield(selector(item))
			})
		},
		size: q.size,
	}
}

// SelectIndexed projects each element of a collection into a new form by
// incorporating the element's index. Returns a query with the result of
// invoking the transform function on each element of the original source.
//
// The first argument to selector represents the zero-based index of that
// element in the source collection. This can be useful if the elements are in a
// known order, and you want to do something with an element at a particular
// index, for example. It can also be useful if you want to retrieve the index
// of one or more elements. The second argument to selector represents the
// element to the process.
func (q Query[T]) SelectIndexed[TResult any](selector func(int, T) TResult) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			index := 0
			q.Iterate(func(item T) bool {
				newItem := selector(index, item)
				index++
				return yield(newItem)
			})
		},
		size: q.size,
	}
}
