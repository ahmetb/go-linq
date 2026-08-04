package linq

import "sort"

type typedSorter[T any] struct {
	items []T
	less  func(T, T) bool
}

func (s typedSorter[T]) Len() int           { return len(s.items) }
func (s typedSorter[T]) Less(i, j int) bool { return s.less(s.items[i], s.items[j]) }
func (s typedSorter[T]) Swap(i, j int)      { s.items[i], s.items[j] = s.items[j], s.items[i] }

func (q query[T]) Sort(less func(T, T) bool) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		items := q.toSlice()
		sort.Sort(typedSorter[T]{items: items, less: less})

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}}
}
