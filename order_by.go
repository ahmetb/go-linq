package linq

import (
	"cmp"
	"slices"
)

type typedOrder[T any] struct {
	compare func(T, T) int
}

// OrderedQuery is a typed query with one or more stable sort keys.
type OrderedQuery[T any] struct {
	Query[T]
	original Query[T]
	orders   []typedOrder[T]
}

func (q Query[T]) OrderBy[K cmp.Ordered](selector func(T) K) OrderedQuery[T] {
	return newOrderedQuery(q, []typedOrder[T]{newTypedOrder(selector)})
}

func (q Query[T]) OrderByDescending[K cmp.Ordered](selector func(T) K) OrderedQuery[T] {
	return newOrderedQuery(q, []typedOrder[T]{newTypedOrderDescending(selector)})
}

func (q OrderedQuery[T]) ThenBy[K cmp.Ordered](selector func(T) K) OrderedQuery[T] {
	orders := append(slices.Clone(q.orders), newTypedOrder(selector))
	return newOrderedQuery(q.original, orders)
}

func (q OrderedQuery[T]) ThenByDescending[K cmp.Ordered](selector func(T) K) OrderedQuery[T] {
	orders := append(slices.Clone(q.orders), newTypedOrderDescending(selector))
	return newOrderedQuery(q.original, orders)
}

func newTypedOrder[T any, K cmp.Ordered](selector func(T) K) typedOrder[T] {
	return typedOrder[T]{
		compare: func(left, right T) int {
			return cmp.Compare(selector(left), selector(right))
		},
	}
}

func newTypedOrderDescending[T any, K cmp.Ordered](selector func(T) K) typedOrder[T] {
	return typedOrder[T]{
		compare: func(left, right T) int {
			return -cmp.Compare(selector(left), selector(right))
		},
	}
}

func newOrderedQuery[T any](original Query[T], orders []typedOrder[T]) OrderedQuery[T] {
	ordered := OrderedQuery[T]{original: original, orders: orders}
	ordered.Query = Query[T]{
		iterate: func(yield func(T) bool) {
			items := original.toSlice()
			slices.SortStableFunc(items, func(left, right T) int {
				for _, order := range orders {
					if result := order.compare(left, right); result != 0 {
						return result
					}
				}
				return 0
			})
			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		},
	}
	return ordered
}
