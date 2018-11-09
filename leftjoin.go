package linq

// LeftJoin correlates the elements of two collection based on matching keys.
//
// A join refers to the operation of correlating the elements of two sources of
// information based on a common key. LeftJoin brings the two information sources
// and the keys by which they are matched together in one method call. This
// differs from the use of Join, as it will also bring elements from outer
// that have no match in inner collection
//
// LeftJoin preserves the order of the elements of outer collection, and for each of
// these elements, the order of the matching elements of inner.
func (q Query) LeftJoin(inner Query,
	outerKeySelector func(interface{}) interface{},
	innerKeySelector func(interface{}) interface{},
	resultSelector func(outer interface{}, inner interface{}) interface{},
	resultSelectorLeftOnly func(outer interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			outernext := q.Iterate()
			innernext := inner.Iterate()

			innerLookup := make(map[interface{}][]interface{})
			for innerItem, ok := innernext(); ok; innerItem, ok = innernext() {
				innerKey := innerKeySelector(innerItem)
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
			}

			var outerItem interface{}
			var innerGroup []interface{}
			innerLen, innerIndex := 0, 0

			return func() (item interface{}, ok bool) {
				if innerIndex >= innerLen {
					has := false
					outerItem, ok = outernext()
					if !ok {
						return nil, false
					}

					innerGroup, has = innerLookup[outerKeySelector(outerItem)]
					if !has {
						return resultSelectorLeftOnly(outerItem), true
					}
					innerLen = len(innerGroup)
					innerIndex = 0
				}

				item = resultSelector(outerItem, innerGroup[innerIndex])
				innerIndex++
				return item, true
			}
		},
	}
}
