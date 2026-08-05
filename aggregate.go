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
