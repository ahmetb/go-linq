package linq

// ZipT is the typed version of Zip.
//
// NOTE: Zip method has better performance than ZipT
//
// resultSelectorFn is of a type "func(TFirst,TSecond)TResult"
//
func (q Query) ZipT(q2 Query, resultSelectorFn interface{}) Query {

	resultSelectorFunc, ok := resultSelectorFn.(func(interface{}, interface{}) interface{})
	if !ok {
		resultSelectorGenericFunc, err := newGenericFunc(
			"ZipT", "resultSelectorFn", resultSelectorFn,
			simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
		)
		if err != nil {
			panic(err)
		}

		resultSelectorFunc = func(item1 interface{}, item2 interface{}) interface{} {
			return resultSelectorGenericFunc.Call(item1, item2)
		}
	}
	return Query{
		Iterate: func() Iterator {
			next1 := q.Iterate()
			next2 := q2.Iterate()

			return func() (item interface{}, ok bool) {
				item1, ok1 := next1()
				item2, ok2 := next2()

				if ok1 && ok2 {
					return resultSelectorFunc(item1, item2), true
				}

				return nil, false
			}
		},
	}
}
