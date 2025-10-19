package linq

// Union produces the set union of two collections.
//
// This method excludes duplicates from the return set. This is different
// behavior to the Concat method, which returns all the elements in the input
// collection, including duplicates.
func (q Query) Union(q2 Query) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			set := make(map[any]struct{})
			stopped := false

			q.Iterate(func(item any) bool {
				if _, seen := set[item]; !seen {
					set[item] = struct{}{}
					if !yield(item) {
						stopped = true
						return false
					}
				}
				return true
			})

			if stopped {
				return
			}

			q2.Iterate(func(item any) bool {
				if _, seen := set[item]; !seen {
					set[item] = struct{}{}
					if !yield(item) {
						return false
					}
				}
				return true
			})
		},
	}
}
