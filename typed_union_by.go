package linq

func (q query[T]) UnionBy[K comparable](other query[T], selector func(T) K) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		seen := make(map[K]struct{})
		continuing := true
		emit := func(value T) bool {
			key := selector(value)
			if _, found := seen[key]; found {
				return true
			}
			seen[key] = struct{}{}
			continuing = yield(value)
			return continuing
		}

		q.iterate(emit)
		if continuing {
			other.iterate(emit)
		}
	}}
}
