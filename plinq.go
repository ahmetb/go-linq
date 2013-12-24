package linq

// ParallelQuery is the type returned from functions executing in parallel.
// To transform a Query into ParallelQuery, use AsParallel() and use
// AsSequential() to do vice versa.
type ParallelQuery struct {
	values  []T
	ordered bool
	err     error
}

type parallelBinaryResult struct {
	ok    bool
	err   error
	index int
}

type parallelValueResult struct {
	val   T
	err   error
	index int
}

// Results evaluates the query and returns the results as T slice.
// An error occurred in during evaluation of the query will be returned.
func (q ParallelQuery) Results() ([]T, error) {
	return q.values, q.err
}

// AsSequential returns a Query from the same source and the query functions
// can be executed in serial for each element of the source sequence.
// This is for undoing AsParallel().
func (q ParallelQuery) AsSequential() Query {
	return Query{values: q.values, err: q.err}
}

// copyMeta copies all fields of ParallelQuery except 'values' into a new
// instance. This should be used for retaining options e.g. 'ordered'.
func (q ParallelQuery) copyMeta() ParallelQuery {
	return ParallelQuery{err: q.err,
		ordered: q.ordered}
}

// copyMetaWithValues copies all fields of ParallelQuery. This should be used
// for retaining options e.g. 'ordered' as well as values.
func (q ParallelQuery) copyMetaWithValues() ParallelQuery {
	return ParallelQuery{err: q.err,
		ordered: q.ordered,
		values:  q.values}
}

// AsOrdered makes the parallel queries to preserve original order. By default,
// parallel queries do not preserve the order and process the parallel
// executions in first-come-first-served fashion.
//
// Not applicable for all query methods and comes
// with a performance penalty in some queries, please refer to
// http://msdn.microsoft.com/en-us/library/dd460677(v=vs.110).aspx .
func (q ParallelQuery) AsOrdered() (p ParallelQuery) {
	p = q.copyMetaWithValues()
	p.ordered = true
	return
}

// AsUnordered undoes the effect of AsOrdered() and do not enforce parallel
// query to preserve the original order.
//
// See AsOrdered() for remarks.
func (q ParallelQuery) AsUnordered() (p ParallelQuery) {
	p = q.copyMetaWithValues()
	p.ordered = true
	return
}

// Where filters a sequence of values by running given predicate function
// in parallel for each element.
//
// This function will take elements of the source (or results of previous query)
// as interface[] so it should make type assertion to work on the types.
// Returns a query with elements satisfy the condition.
//
// If any of the parallel executions return with an error, this function
// immediately returns with the error.
//
// If you would like to preserve order from the original sequence, use
// AsOrdered() on the query beforehand.
func (q ParallelQuery) Where(f func(T) (bool, error)) (r ParallelQuery) {
	r = q.copyMeta()
	if r.err != nil {
		return r
	}
	if f == nil {
		r.err = ErrNilFunc
		return
	}

	count := len(q.values)
	ch := make(chan *parallelBinaryResult)
	for i := 0; i < count; i++ {
		go func(ind int, f func(T) (bool, error), in T) {
			out := parallelBinaryResult{index: ind}
			ok, err := f(in)
			if err != nil {
				out.err = err
			} else {
				out.ok = ok
			}
			ch <- &out
		}(i, f, q.values[i])
	}

	tmp := make([]T, count)
	take := make([]bool, count)

	for j := 0; j < count; j++ {
		out := <-ch
		if out.err != nil {
			r.err = out.err
			return
		}
		if out.ok {
			origI := out.index
			val := q.values[origI]
			if r.ordered {
				tmp[origI] = val
				take[origI] = true
			} else {
				r.values = append(r.values, val)
			}
		}
	}

	if r.ordered {
		// iterate over the flag slice to take marked elements
		for i, v := range tmp {
			if take[i] {
				r.values = append(r.values, v)
			}
		}
	}
	return
}

// Select projects each element of a sequence into a new form by running
// the given transform function in parallel for each element.
// Returns a query with the return values of invoking the transform function
// on each element of original source.
func (q ParallelQuery) Select(f func(T) (T, error)) (r ParallelQuery) {
	r = q.copyMeta()
	if r.err != nil {
		return r
	}
	if f == nil {
		r.err = ErrNilFunc
		return
	}

	ch := make(chan *parallelValueResult)
	r.values = make([]T, len(q.values))
	for i, v := range q.values {
		go func(ind int, f func(T) (T, error), in T) {
			out := parallelValueResult{index: ind}
			val, err := f(in)
			if err != nil {
				out.err = err
			} else {
				out.val = val
			}
			ch <- &out
		}(i, f, v)
	}

	for i := 0; i < len(q.values); i++ {
		out := <-ch
		if out.err != nil {
			r.err = out.err
			return
		}
		r.values[out.index] = out.val
	}
	return
}

// AnyWith determines whether the query source contains any elements satisfying
// the provided predicate function.
func (q ParallelQuery) AnyWith(f func(T) (bool, error)) (exists bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	ch := make(chan parallelBinaryResult)
	for _, v := range q.values {
		go func(f func(T) (bool, error), value T) {
			out := parallelBinaryResult{}
			ok, e := f(value)
			out.ok = ok
			out.err = e
			ch <- out
		}(f, v)
	}

	for i := 0; i < len(q.values); i++ {
		out := <-ch
		if out.err != nil {
			err = out.err
			return
		}
		if out.ok {
			exists = true
			return
		}
	}
	return
}

// All determines whether all elements of the query source satisfy the provided
// predicate function by executing the function for each element in parallel.
//
// Returns early if one element does not meet the conditions provided.
func (q ParallelQuery) All(f func(T) (bool, error)) (all bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	ch := make(chan parallelBinaryResult)
	for v := range q.values {
		go func(f func(T) (bool, error), value T) {
			ok, e := f(value)
			ch <- parallelBinaryResult{ok: ok, err: e}
		}(f, v)
	}

	for i := 0; i < len(q.values); i++ {
		out := <-ch
		if out.err != nil {
			err = out.err
			return
		}
		if !out.ok {
			return false, nil
		}
	}
	return true, nil
}
