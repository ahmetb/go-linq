package linq

import (
	"cmp"
	"slices"
)

type typedOrder[T any] struct {
	compare func(T, T) int
}

type orderedQuery[T any] struct {
	query[T]
	original query[T]
	orders   []typedOrder[T]
}

func (q query[T]) OrderBy[K cmp.Ordered](selector func(T) K) orderedQuery[T] {
	return newOrderedQuery(q, []typedOrder[T]{newTypedOrder(selector)})
}

func (q query[T]) OrderByDescending[K cmp.Ordered](selector func(T) K) orderedQuery[T] {
	return newOrderedQuery(q, []typedOrder[T]{newTypedOrderDescending(selector)})
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

func newOrderedQuery[T any](original query[T], orders []typedOrder[T]) orderedQuery[T] {
	ordered := orderedQuery[T]{original: original, orders: orders}
	ordered.query = query[T]{
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
