package linq

// Aggregate applies an accumulator function over a sequence.
//
// Aggregate method makes it simple to perform a calculation over a sequence of values.
// This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both
// the element from the sequence and an aggregated value (as the first argument to f()).
// The first element of source is used as the initial aggregate value.
// The result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (q Query) Aggregate(
	f func(interface{}, interface{}) interface{}) interface{} {
	next := q.Iterate()

	result, any := next()
	if !any {
		return nil
	}

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return result
}

// AggregateWithSeed applies an accumulator function over a sequence.
// The specified seed value is used as the initial accumulator value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of values.
// This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both
// the element from the sequence and an aggregated value (as the first argument to f()).
// The value of the seed parameter is used as the initial aggregate value.
// The result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (q Query) AggregateWithSeed(
	seed interface{},
	f func(interface{}, interface{}) interface{},
) interface{} {

	next := q.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = f(result, current)
	}

	return result
}
