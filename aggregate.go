package linq

// Aggregate combines the sequence from left to right, using its first value as the seed.
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

// AggregateWithSeed combines the sequence from left to right starting with seed.
func (q Query[T]) AggregateWithSeed[A any](seed A, accumulator func(A, T) A) A {
	result := seed
	for value := range q.iterate {
		result = accumulator(result, value)
	}
	return result
}

// AggregateWithSeedBy aggregates from seed and projects the final accumulator.
func (q Query[T]) AggregateWithSeedBy[A, R any](
	seed A,
	accumulator func(A, T) A,
	resultSelector func(A) R,
) R {
	return resultSelector(q.AggregateWithSeed(seed, accumulator))
}
