package linq

// Union produces the set union of two collections.
//
// This method excludes duplicates from the return set. This is different
// behavior to the Concat method, which returns all the elements in the input
// collection including duplicates.
func (q Query) Union(q2 Query) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			next2 := q2.Iterate()

			set := make(map[interface{}]bool)
			use1 := true

			return func() (item interface{}, ok bool) {
				if use1 {
					for item, ok = next(); ok; item, ok = next() {
						if _, has := set[item]; !has {
							set[item] = true
							return
						}
					}

					use1 = false
				}

				for item, ok = next2(); ok; item, ok = next2() {
					if _, has := set[item]; !has {
						set[item] = true
						return
					}
				}

				return
			}
		},
	}
}
