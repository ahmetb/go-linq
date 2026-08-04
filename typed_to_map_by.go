package linq

func (q query[T]) ToMapBy[K comparable, V any](
	keySelector func(T) K,
	valueSelector func(T) V,
) map[K]V {
	result := make(map[K]V)
	for value := range q.iterate {
		result[keySelector(value)] = valueSelector(value)
	}

	return result
}
