package linq

// DefaultIfEmpty returns the elements of the specified sequence
// if the sequence is empty.
func (q Query) DefaultIfEmpty(defaultValue interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			state := 1

			return func() (item interface{}, ok bool) {
				switch state {
				case 1:
					item, ok = next()
					if ok {
						state = 2
					} else {
						item = defaultValue
						ok = true
						state = -1
					}
					return
				case 2:
					for item, ok = next(); ok; item, ok = next() {
						return
					}
					return
				}
				return
			}
		},
	}
}

func (q QueryG[T]) DefaultIfEmpty(defaultValue T) QueryG[T] {
	return QueryG[T]{
		Iterate: func() IteratorG[T] {
			next := q.Iterate()
			state := 1

			return func() (item T, ok bool) {
				switch state {
				case 1:
					item, ok = next()
					if ok {
						state = 2
					} else {
						item = defaultValue
						ok = true
						state = -1
					}
					return
				case 2:
					for item, ok = next(); ok; item, ok = next() {
						return
					}
					return
				}
				return
			}
		},
	}
}
