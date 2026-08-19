package linq

// Aggregate applies an accumulator function over a sequence.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in a source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The first element of the source is used as the initial aggregate value. The
// result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f() and a boolean reporting whether
// the sequence was non-empty.
func (q Query[T]) Aggregate(f func(accumulator, item T) T) (T, bool) {
	var result T
	first := true

	q.Iterate(func(current T) bool {
		if first {
			result = current
			first = false
		} else {
			result = f(result, current)
		}
		return true
	})

	return result, !first
}

// AggregateWithSeed applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in a source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The value of the seed parameter is used as the initial aggregate value.
// The result of f() replaces the previous aggregated value.
//
// AggregateWithSeed is a generic method: the accumulator type TAccumulate is inferred
// from the seed value and may differ from the element type T.
//
// Aggregate returns the final result of f().
func (q Query[T]) AggregateWithSeed[TAccumulate any](seed TAccumulate,
	f func(accumulator TAccumulate, item T) TAccumulate) TAccumulate {
	result := seed

	q.Iterate(func(current T) bool {
		result = f(result, current)
		return true
	})

	return result
}

// AggregateWithSeedBy applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value, and the
// specified function is used to select the result value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source.
// Each time func is called, Aggregate passes both the element from the sequence
// and an aggregated value (as the first argument to func). The value of the
// seed parameter is used as the initial aggregate value. The result of func
// replaces the previous aggregated value.
//
// AggregateWithSeedBy is a generic method: the accumulator type TAccumulate is inferred
// from the seed value, and the result type TResult from the resultSelector function.
//
// The final result of func is passed to resultSelector to obtain the final
// result of Aggregate.
func (q Query[T]) AggregateWithSeedBy[TAccumulate, TResult any](seed TAccumulate,
	f func(accumulator TAccumulate, item T) TAccumulate,
	resultSelector func(TAccumulate) TResult) TResult {

	result := seed

	q.Iterate(func(current T) bool {
		result = f(result, current)
		return true
	})

	return resultSelector(result)
}

// CountBy returns the number of elements for each selected key. Results are
// ordered by the key's first appearance in the source.
//
// CountBy is a generic method: the key type TKey is inferred from the
// keySelector function. The key type TKey must be comparable.
func (q Query[T]) CountBy[TKey comparable](keySelector func(T) TKey) Query[KeyValue[TKey, int]] {
	// Keep this loop specialized: delegating to aggregateBy adds a closure call
	// per element and two allocations per query.
	return Query[KeyValue[TKey, int]]{
		Iterate: func(yield func(KeyValue[TKey, int]) bool) {
			index := make(map[TKey]int)
			var counts []KeyValue[TKey, int]

			q.Iterate(func(item T) bool {
				key := keySelector(item)
				i, exists := index[key]
				if !exists {
					i = len(counts)
					index[key] = i
					counts = append(counts, KeyValue[TKey, int]{Key: key})
				}
				counts[i].Value++
				return true
			})

			for _, count := range counts {
				if !yield(count) {
					return
				}
			}
		},
	}
}

// AggregateBy applies an accumulator to elements grouped by a selected key.
// Each key starts with seed, and results are ordered by the key's first
// appearance in the source.
//
// AggregateBy is a generic method: the key type TKey is inferred from the
// keySelector function and the accumulator type TAccumulate from the seed
// value. The key type TKey must be comparable.
func (q Query[T]) AggregateBy[TKey comparable, TAccumulate any](keySelector func(T) TKey,
	seed TAccumulate,
	f func(accumulator TAccumulate, item T) TAccumulate) Query[KeyValue[TKey, TAccumulate]] {
	return aggregateBy(q, keySelector, func(TKey) TAccumulate { return seed }, f)
}

// AggregateByWithSeedSelector applies an accumulator to elements grouped by a
// selected key. seedSelector supplies the initial accumulator once per key,
// and results are ordered by the key's first appearance in the source.
//
// AggregateByWithSeedSelector is a generic method: the key type TKey is
// inferred from the keySelector function and the accumulator type TAccumulate
// from the seedSelector function. The key type TKey must be comparable.
func (q Query[T]) AggregateByWithSeedSelector[TKey comparable, TAccumulate any](
	keySelector func(T) TKey,
	seedSelector func(TKey) TAccumulate,
	f func(accumulator TAccumulate, item T) TAccumulate) Query[KeyValue[TKey, TAccumulate]] {
	return aggregateBy(q, keySelector, seedSelector, f)
}

func aggregateBy[T any, TKey comparable, TAccumulate any](q Query[T],
	keySelector func(T) TKey,
	seedSelector func(TKey) TAccumulate,
	f func(TAccumulate, T) TAccumulate) Query[KeyValue[TKey, TAccumulate]] {
	return Query[KeyValue[TKey, TAccumulate]]{
		Iterate: func(yield func(KeyValue[TKey, TAccumulate]) bool) {
			index := make(map[TKey]int)
			var aggregates []KeyValue[TKey, TAccumulate]

			q.Iterate(func(item T) bool {
				key := keySelector(item)
				i, exists := index[key]
				if !exists {
					i = len(aggregates)
					index[key] = i
					aggregates = append(aggregates, KeyValue[TKey, TAccumulate]{
						Key:   key,
						Value: seedSelector(key),
					})
				}
				aggregates[i].Value = f(aggregates[i].Value, item)
				return true
			})

			for _, aggregate := range aggregates {
				if !yield(aggregate) {
					return
				}
			}
		},
	}
}
