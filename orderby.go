package linq

import "sort"

type order struct {
	selector func(interface{}) interface{}
	compare  comparer
	desc     bool
}

// OrderedQuery is the type returned from OrderBy, OrderByDescending ThenBy and
// ThenByDescending functions.
type OrderedQuery struct {
	Query
	original Query
	orders   []order
}

// OrderBy sorts the elements of a collection in ascending order. Elements are
// sorted according to a key.
func (q Query) OrderBy(selector func(interface{}) interface{}) OrderedQuery {
	return OrderedQuery{
		orders:   []order{{selector: selector}},
		original: q,
		Query: Query{
			Iterate: func() Iterator {
				items := q.sort([]order{{selector: selector}})
				len := len(items)
				index := 0

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}
}

// OrderByT is the typed version of OrderBy.
//
//   - selectorFn is of type "func(TSource) TKey"
//
// NOTE: OrderBy has better performance than OrderByT.
func (q Query) OrderByT(selectorFn interface{}) OrderedQuery {
	selectorGenericFunc, err := newGenericFunc(
		"OrderByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.OrderBy(selectorFunc)
}

// OrderByDescending sorts the elements of a collection in descending order.
// Elements are sorted according to a key.
func (q Query) OrderByDescending(selector func(interface{}) interface{}) OrderedQuery {
	return OrderedQuery{
		orders:   []order{{selector: selector, desc: true}},
		original: q,
		Query: Query{
			Iterate: func() Iterator {
				items := q.sort([]order{{selector: selector, desc: true}})
				len := len(items)
				index := 0

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}
}

// OrderByDescendingT is the typed version of OrderByDescending.
//   - selectorFn is of type "func(TSource) TKey"
// NOTE: OrderByDescending has better performance than OrderByDescendingT.
func (q Query) OrderByDescendingT(selectorFn interface{}) OrderedQuery {
	selectorGenericFunc, err := newGenericFunc(
		"OrderByDescendingT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return q.OrderByDescending(selectorFunc)
}

// ThenBy performs a subsequent ordering of the elements in a collection in
// ascending order. This method enables you to specify multiple sort criteria by
// applying any number of ThenBy or ThenByDescending methods.
func (oq OrderedQuery) ThenBy(
	selector func(interface{}) interface{}) OrderedQuery {
	return OrderedQuery{
		orders:   append(oq.orders, order{selector: selector}),
		original: oq.original,
		Query: Query{
			Iterate: func() Iterator {
				items := oq.original.sort(append(oq.orders, order{selector: selector}))
				len := len(items)
				index := 0

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}
}

// ThenByT is the typed version of ThenBy.
//   - selectorFn is of type "func(TSource) TKey"
// NOTE: ThenBy has better performance than ThenByT.
func (oq OrderedQuery) ThenByT(selectorFn interface{}) OrderedQuery {
	selectorGenericFunc, err := newGenericFunc(
		"ThenByT", "selectorFn", selectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	selectorFunc := func(item interface{}) interface{} {
		return selectorGenericFunc.Call(item)
	}

	return oq.ThenBy(selectorFunc)
}

// ThenByDescending performs a subsequent ordering of the elements in a
// collection in descending order. This method enables you to specify multiple
// sort criteria by applying any number of ThenBy or ThenByDescending methods.
func (oq OrderedQuery) ThenByDescending(selector func(interface{}) interface{}) OrderedQuery {
	return OrderedQuery{
		orders:   append(oq.orders, order{selector: selector, desc: true}),
		original: oq.original,
		Query: Query{
			Iterate: func() Iterator {
				items := oq.original.sort(append(oq.orders, order{selector: selector, desc: true}))
				len := len(items)
				index := 0

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}
}

// ThenByDescendingT is the typed version of ThenByDescending.
//   - selectorFn is of type "func(TSource) TKey"
// NOTE: ThenByDescending has better performance than ThenByDescendingT.
func (oq OrderedQuery) ThenByDescendingT(selectorFn interface{}) OrderedQuery {
	selectorFunc, ok := selectorFn.(func(interface{}) interface{})
	if !ok {
		selectorGenericFunc, err := newGenericFunc(
			"ThenByDescending", "selectorFn", selectorFn,
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
		)
		if err != nil {
			panic(err)
		}

		selectorFunc = func(item interface{}) interface{} {
			return selectorGenericFunc.Call(item)
		}
	}
	return oq.ThenByDescending(selectorFunc)
}

// Sort returns a new query by sorting elements with provided less function in
// ascending order. The comparer function should return true if the parameter i
// is less than j. While this method is uglier than chaining OrderBy,
// OrderByDescending, ThenBy and ThenByDescending methods, it's performance is
// much better.
func (q Query) Sort(less func(i, j interface{}) bool) Query {
	return Query{
		Iterate: func() Iterator {
			items := q.lessSort(less)
			len := len(items)
			index := 0

			return func() (item interface{}, ok bool) {
				ok = index < len
				if ok {
					item = items[index]
					index++
				}

				return
			}
		},
	}
}

// SortT is the typed version of Sort.
//   - lessFn is of type "func(TSource,TSource) bool"
// NOTE: Sort has better performance than SortT.
func (q Query) SortT(lessFn interface{}) Query {
	lessGenericFunc, err := newGenericFunc(
		"SortT", "lessFn", lessFn,
		simpleParamValidator(newElemTypeSlice(new(genericType), new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	lessFunc := func(i, j interface{}) bool {
		return lessGenericFunc.Call(i, j).(bool)
	}

	return q.Sort(lessFunc)
}

type sorter struct {
	items []interface{}
	less  func(i, j interface{}) bool
}

func (s sorter) Len() int {
	return len(s.items)
}

func (s sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s sorter) Less(i, j int) bool {
	return s.less(s.items[i], s.items[j])
}

func (q Query) sort(orders []order) (r []interface{}) {
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		r = append(r, item)
	}

	if len(r) == 0 {
		return
	}

	for i, j := range orders {
		orders[i].compare = getComparer(j.selector(r[0]))
	}

	s := sorter{
		items: r,
		less: func(i, j interface{}) bool {
			for _, order := range orders {
				x, y := order.selector(i), order.selector(j)
				switch order.compare(x, y) {
				case 0:
					continue
				case -1:
					return !order.desc
				default:
					return order.desc
				}
			}

			return false
		}}

	sort.Sort(s)
	return
}

func (q Query) lessSort(less func(i, j interface{}) bool) (r []interface{}) {
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		r = append(r, item)
	}

	s := sorter{items: r, less: less}

	sort.Sort(s)
	return
}
