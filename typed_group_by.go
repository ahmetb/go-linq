package linq

type group[K comparable, E any] struct {
	Key   K
	Group []E
}

func (q Query[T]) GroupBy[K comparable, E any](
	keySelector func(T) K,
	elementSelector func(T) E,
) Query[group[K, E]] {
	return Query[group[K, E]]{
		iterate: func(yield func(group[K, E]) bool) {
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
				if !yield(group[K, E]{Key: key, Group: groups[key]}) {
					return
				}
			}
		},
	}
}
