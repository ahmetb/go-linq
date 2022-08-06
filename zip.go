package linq

// Zip applies a specified function to the corresponding elements of two
// collections, producing a collection of the results.
//
// The method steps through the two input collections, applying function
// resultSelector to corresponding elements of the two collections. The method
// returns a collection of the values that are returned by resultSelector. If
// the input collections do not have the same number of elements, the method
// combines elements until it reaches the end of one of the collections. For
// example, if one collection has three elements and the other one has four, the
// result collection has only three elements.
func (q Query) Zip(q2 Query,
	resultSelector func(interface{}, interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			next1 := q.Iterate()
			next2 := q2.Iterate()

			return func() (item interface{}, ok bool) {
				item1, ok1 := next1()
				item2, ok2 := next2()

				if ok1 && ok2 {
					return resultSelector(item1, item2), true
				}

				return nil, false
			}
		},
	}
}

func (q QueryG[T]) Zip(selector Zipper[T]) interface{} {
	return selector.Map(q)
}

type Zipper[TIn1 any] interface {
	Map(q QueryG[TIn1]) interface{}
}

func ZipWith[TIn1, TIn2, TOut any](q2 QueryG[TIn2], resultSelector func(TIn1, TIn2) TOut) Zipper[TIn1] {
	return zipSelector[TIn1, TIn2, TOut]{
		resultSelector: resultSelector,
		q2:             q2,
	}
}

type zipSelector[TIn1, TIn2, TOut any] struct {
	resultSelector func(TIn1, TIn2) TOut
	q2             QueryG[TIn2]
}

var z Zipper[int] = zipSelector[int, int, int]{}

func (z zipSelector[TIn1, TIn2, TOut]) Map(q QueryG[TIn1]) interface{} {
	return ZipG(q, z.q2, z.resultSelector)
}

func ZipG[TIn1, TIn2, TOut any](q QueryG[TIn1], q2 QueryG[TIn2],
	resultSelector func(TIn1, TIn2) TOut) QueryG[TOut] {
	return QueryG[TOut]{
		Iterate: func() IteratorG[TOut] {
			next1 := q.Iterate()
			next2 := q2.Iterate()

			return func() (item TOut, ok bool) {
				item1, ok1 := next1()
				item2, ok2 := next2()

				if ok1 && ok2 {
					return resultSelector(item1, item2), true
				}

				return *new(TOut), false
			}
		},
	}
}

// ZipT is the typed version of Zip.
//
//   - resultSelectorFn is of type "func(TFirst,TSecond)TResult"
//
// NOTE: Zip has better performance than ZipT.
func (q Query) ZipT(q2 Query,
	resultSelectorFn interface{}) Query {
	resultSelectorGenericFunc, err := newGenericFunc(
		"ZipT", "resultSelectorFn", resultSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	resultSelectorFunc := func(item1 interface{}, item2 interface{}) interface{} {
		return resultSelectorGenericFunc.Call(item1, item2)
	}

	return q.Zip(q2, resultSelectorFunc)
}
