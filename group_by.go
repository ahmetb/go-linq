package linq

// Group contains one key and the source elements associated with it.
type Group[K comparable, E any] struct {
	Key   K
	Group []E
}

func (q Query[T]) GroupBy[K comparable, E any](
	keySelector func(T) K,
	elementSelector func(T) E,
) Query[Group[K, E]] {
	return Query[Group[K, E]]{
		iterate: func(yield func(Group[K, E]) bool) {
			groups := make(map[K][]E)
			var keys []K
			for value := range q.iterate {
				key := keySelector(value)
				if _, ok := groups[key]; !ok {
					keys = append(keys, key)
				}
				groups[key] = append(groups[key], elementSelector(value))
			}

			for _, key := range keys {
				if !yield(Group[K, E]{Key: key, Group: groups[key]}) {
					return
				}
			}
		},
	}
}
