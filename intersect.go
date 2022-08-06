package linq

// Intersect produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
func (q Query) Intersect(q2 Query) Query {
	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			next2 := q2.Iterate()

			set := make(map[interface{}]bool)
			for item, ok := next2(); ok; item, ok = next2() {
				set[item] = true
			}

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if _, has := set[item]; has {
						delete(set, item)
						return
					}
				}

				return
			}
		},
	}
}

// Intersect produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
func (q QueryG[T]) Intersect(q2 QueryG[T]) QueryG[T] {
	return QueryG[T]{
		Iterate: func() IteratorG[T] {
			next := q.Iterate()
			next2 := q2.Iterate()

			set := make(map[interface{}]bool)
			for item, ok := next2(); ok; item, ok = next2() {
				set[item] = true
			}

			return func() (item T, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if _, has := set[item]; has {
						delete(set, item)
						return
					}
				}

				return
			}
		},
	}
}

// IntersectBy produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
//
// IntersectBy invokes a transform function on each element of both collections.
func (q Query) IntersectBy(q2 Query,
	selector func(interface{}) interface{}) Query {

	return Query{
		Iterate: func() Iterator {
			next := q.Iterate()
			next2 := q2.Iterate()

			set := make(map[interface{}]bool)
			for item, ok := next2(); ok; item, ok = next2() {
				s := selector(item)
				set[s] = true
			}

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					s := selector(item)
					if _, has := set[s]; has {
						delete(set, s)
						return
					}
				}

				return
			}
		},
	}
}

type Intersect[T any] interface {
	Map(q1, q2 QueryG[T]) any
}

type intersectBy[TIn, TOut any] struct {
	selector func(TIn) TOut
}

func IntersectSelector[TIn, TOut any](selector func(TIn) TOut) Intersect[TIn] {
	return intersectBy[TIn, TOut]{
		selector: selector,
	}
}

func (i intersectBy[TIn, TOut]) Map(q1, q2 QueryG[TIn]) any {
	return IntersectBy[TIn, TOut](q1, q2, i.selector)
}

func (q QueryG[T]) IntersectBy(q2 QueryG[T], i Intersect[T]) QueryG[T] {
	return i.Map(q, q2).(QueryG[T])
}

func IntersectBy[T1, T2 any](q1, q2 QueryG[T1], selector func(T1) T2) QueryG[T1] {
	return QueryG[T1]{
		Iterate: func() IteratorG[T1] {
			next := q1.Iterate()
			next2 := q2.Iterate()

			set := make(map[interface{}]bool)
			for item, ok := next2(); ok; item, ok = next2() {
				s := selector(item)
				set[s] = true
			}

			return func() (item T1, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					s := selector(item)
					if _, has := set[s]; has {
						delete(set, s)
						return
					}
				}

				return
			}
		},
	}
}

// IntersectByT is the typed version of IntersectBy.
//
//   - selectorFn is of type "func(TSource) TSource"
//
// NOTE: IntersectBy has better performance than IntersectByT.
func (q Query) IntersectByT(q2 Query,
	selectorFn interface{}) Query {
	selectorGenericFunc, err := newGenericFunc(
		"IntersectByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.IntersectBy(q2, selectorFunc)
}
