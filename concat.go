package linq

// Append inserts an item to the end of a collection, so it becomes the last
// item.
func (q Query) Append(item interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			appended := false

			return func() (interface{}, bool) {
				i, ok := next()
				if ok {
					return i, ok
				}

				if !appended {
					appended = true
					return item, true
				}

				return nil, false
			}
		},
	}
}

func (q QueryG[T]) Append(item T) QueryG[T] {
	return QueryG[T]{
		Iterate: func() IteratorG[T] {
			next := q.Iterate()
			appended := false

			return func() (T, bool) {
				i, ok := next()
				if ok {
					return i, ok
				}

				if !appended {
					appended = true
					return item, true
				}

				return *new(T), false
			}
		},
	}
}

// Concat concatenates two collections.
//
// The Concat method differs from the Union method because the Concat method
// returns all the original elements in the input sequences. The Union method
// returns only unique elements.
func (q Query) Concat(q2 Query) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			next2 := q2.Iterate()
			use1 := true

			return func() (item interface{}, ok bool) {
				if use1 {
					item, ok = next()
					if ok {
						return
					}

					use1 = false
				}

				return next2()
			}
		},
	}
}

func (q QueryG[T]) Concat(q2 QueryG[T]) QueryG[T] {
	return QueryG[T]{
		Iterate: func() IteratorG[T] {
			next := q.Iterate()
			next2 := q2.Iterate()
			use1 := true

			return func() (item T, ok bool) {
				if use1 {
					item, ok = next()
					if ok {
						return
					}

					use1 = false
				}

				return next2()
			}
		},
	}
}

// Prepend inserts an item to the beginning of a collection, so it becomes the
// first item.
func (q Query) Prepend(item interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			prepended := false

			return func() (interface{}, bool) {
				if prepended {
					return next()
				}

				prepended = true
				return item, true
			}
		},
	}
}

func (q QueryG[T]) Prepend(item T) QueryG[T] {
	return QueryG[T]{
		Iterate: func() IteratorG[T] {
			next := q.Iterate()
			prepended := false

			return func() (T, bool) {
				if prepended {
					return next()
				}

				prepended = true
				return item, true
			}
		},
	}
}
