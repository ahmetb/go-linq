package linq

func (q Query[T]) Aggregate(accumulator func(T, T) T) T {
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

func (q Query[T]) AggregateWithSeed[A any](seed A, accumulator func(A, T) A) A {
	result := seed
	for value := range q.iterate {
		result = accumulator(result, value)
	}
	return result
}

func (q Query[T]) AggregateWithSeedBy[A, R any](
	seed A,
	accumulator func(A, T) A,
	resultSelector func(A) R,
) R {
	return resultSelector(q.AggregateWithSeed(seed, accumulator))
}
