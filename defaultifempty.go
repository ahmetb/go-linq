package linq

// DefaultIfEmpty returns the elements of the specified sequence
// if the sequence is empty.
func (q Query) DefaultIfEmpty(defaultValue any) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			var yieldedAnyThing bool
			var stopped bool

			q.Iterate(func(item any) bool {
				yieldedAnyThing = true

				if !yield(item) {
					stopped = true
					return false
				}
				return true
			})

			// If the iteration wasn't stopped and the source was empty...
			if !stopped && !yieldedAnyThing {
				// ...yield the default value.
				yield(defaultValue)
			}
		},
	}
}
