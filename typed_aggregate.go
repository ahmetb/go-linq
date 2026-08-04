package linq

func (q query[T]) Aggregate(accumulator func(T, T) T) T {
	var result T
	first := true
	for value := range q.iterate {
		if first {
			result = value
			first = false
			continue
		}
		result = accumulator(result, value)
	}
	return result
}
