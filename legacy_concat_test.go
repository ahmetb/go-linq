package linq

// Append inserts an item to the end of a collection, so it becomes the last
// item.
func (q legacyQuery) Append(item any) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			stopped := false

			q.Iterate(func(originalItem any) bool {
				if !yield(originalItem) {
					stopped = true
					return false
				}
				return true
			})

			if !stopped {
				yield(item)
			}
		},
	}
}

// Concat concatenates two collections.
//
// The Concat method differs from the Union method because the Concat method
// returns all the original elements in the input sequences. The Union method
// returns only unique elements.
func (q legacyQuery) Concat(q2 legacyQuery) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			stopped := false

			q.Iterate(func(item any) bool {
				if !yield(item) {
					stopped = true
					return false
				}
				return true
			})

			if !stopped {
				q2.Iterate(func(item any) bool {
					return yield(item)
				})
			}
		},
	}
}

// Prepend inserts an item to the beginning of a collection, so it becomes the
// first item.
func (q legacyQuery) Prepend(item any) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			if !yield(item) {
				return
			}

			q.Iterate(func(item any) bool {
				return yield(item)
			})
		},
	}
}
