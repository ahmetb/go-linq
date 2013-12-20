// Package linq provides methods for querying and manipulating slices and
// collections.
package linq

import (
	"errors"
	"sort"
)

// T is an alias for interface{} to make things shorter. Whatever you pass
// to linq functions (e.g. From, Union, Intersect, Except) should be a struct
// of this type, []T.
type T interface{}

// Queryable is the type returned from query functions. To evaluante
// get the results of the query, use Results().
type Queryable struct {
	values []T
	err    error
}

type sortableQueryable struct {
	values []T
	less   func(this, that T) bool
}

func (q sortableQueryable) Len() int           { return len(q.values) }
func (q sortableQueryable) Swap(i, j int)      { q.values[i], q.values[j] = q.values[j], q.values[i] }
func (q sortableQueryable) Less(i, j int) bool { return q.less(q.values[i], q.values[j]) }

var (
	ErrNilFunc       = errors.New("linq: passed evaluation function is nil")                                           // a predicate, selector or comparer is nil
	ErrNilInput      = errors.New("linq: nil sequence passed as input to function")                                    // nil value of []T is passed
	ErrNoElement     = errors.New("linq: element satisfying the conditions does not exist")                            // strictly element requesting methods are called and element is not found
	ErrEmptySequence = errors.New("linq: empty sequence, operation requires non-empty results sequence")               // requested operation is invalid on empty sequences
	ErrNegativeParam = errors.New("linq: parameter cannot be negative")                                                // negative value passed to an index parameter
	ErrNan           = errors.New("linq: sequence contains an element of non-numeric types")                           // sequence has invalid elements that method cannot assert into one of builtin numeric types
	ErrTypeMismatch  = errors.New("linq: sequence contains element(s) with type different than requested type or nil") // sequence elements or nil of different type than function can work with
	ErrNotSingle     = errors.New("linq: sequence contains more than one element matching the given predicate found")  // sequence contains more than one elements satisfy given predicate func
)

// From initializes a linq query with passed slice as the source.
// The slice has to be of type []T. This is a language limitation.
func From(input []T) Queryable {
	var _err error
	if input == nil {
		_err = ErrNilInput
	}
	return Queryable{
		values: input,
		err:    _err}
}

// Results evaluates the query and returns the results as T slice.
// An error occurred in during evaluation of the query will be returned.
func (q Queryable) Results() ([]T, error) {
	return q.values, q.err
}

// Where filters a sequence of values based on a predicate function. This
// function will take elements of the source (or results of previous query)
// as interface[] so it should make type assertion to work on the types.
// Returns a query with elements satisfy the condition.
func (q Queryable) Where(f func(T) (bool, error)) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return r
	}
	if f == nil {
		r.err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		ok, err := f(i)
		if err != nil {
			r.err = err
			return r
		}
		if ok {
			r.values = append(r.values, i)
		}
	}
	return r
}

// Select projects each element of a sequence into a new form.
// Returns a query with the result of invoking the transform function
// on each element of original source.
func (q Queryable) Select(f func(T) (T, error)) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return r
	}
	if f == nil {
		r.err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		val, err := f(i)
		if err != nil {
			r.err = err
			return r
		}
		r.values = append(r.values, val)
	}
	return
}

// Distinct returns distinct elements from the provided source using default
// equality comparer, ==. This is a set operation and returns an unordered
// sequence.
func (q Queryable) Distinct() (r Queryable) {
	return q.distinct(nil)
}

// DistinctBy returns distinct elements from the provided source using the
// provided equality comparer. This is a set operation and returns an unordered
// sequence. Number of calls to f will be at most N^2 (all elements are
// distinct) and at best N (all elements are the same).
func (q Queryable) DistinctBy(f func(T, T) (bool, error)) (r Queryable) {
	if f == nil {
		r.err = ErrNilFunc
		return
	}
	return q.distinct(f)
}

// distinct returns distinct elements from the provided source using default
// equality comparer (==) or a custom equality comparer function. Complexity
// is O(N).
func (q Queryable) distinct(f func(T, T) (bool, error)) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return r
	}

	if f == nil {
		// basic equality comparison using dict
		dict := make(map[T]bool)
		for _, v := range q.values {
			if _, ok := dict[v]; !ok {
				dict[v] = true
			}
		}
		res := make([]T, len(dict))
		i := 0
		for key, _ := range dict {
			res[i] = key
			i++
		}
		r.values = res
	} else {
		// use equality comparer and bool flags for each item
		// here we check all a[i]==a[j] i<j, practically worst case
		// for this is O(N^2) where all elements are different and best case
		// is O(N) where all elements are the same
		// pick lefthand side value of the comparison in the result
		l := len(q.values)
		results := make([]T, 0)
		included := make([]bool, l)
		for i := 0; i < l; i++ {
			if included[i] {
				continue
			}
			for j := i + 1; j < l; j++ {
				equals, err := f(q.values[i], q.values[j])
				if err != nil {
					r.err = err
					return
				}
				if equals {
					included[j] = true // don't include righthand side value
				}
			}
			results = append(results, q.values[i])
		}
		r.values = results
	}
	return
}

// Union returns set union of the source sequence and the provided
// input slice using default equality comparer. This is a set operation and
// returns an unordered sequence.
func (q Queryable) Union(in []T) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if in == nil {
		r.err = ErrNilInput
		return
	}
	var set map[T]bool = make(map[T]bool)

	for _, v := range q.values {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	for _, v := range in {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	r.values = make([]T, len(set))
	i := 0
	for k, _ := range set {
		r.values[i] = k
		i++
	}
	return
}

// Intersect returns set intersection of the source sequence and the
// provided input slice using default equality comparer. This is a set
// operation and may return an unordered sequence.
func (q Queryable) Intersect(in []T) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if in == nil {
		r.err = ErrNilInput
		return
	}
	var set map[T]bool = make(map[T]bool)
	var intersection map[T]bool = make(map[T]bool)

	for _, v := range q.values {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	for _, v := range in {
		if _, ok := set[v]; ok {
			delete(set, v)
			if _, added := intersection[v]; !added {
				intersection[v] = true
			}
		}
	}
	r.values = make([]T, len(intersection))
	i := 0
	for k, _ := range intersection {
		r.values[i] = k
		i++
	}
	return
}

// Except returns set difference of the source sequence and the
// provided input slice using default equality comparer. This is a set
// operation and returns an unordered sequence.
func (q Queryable) Except(in []T) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if in == nil {
		r.err = ErrNilInput
		return
	}
	var set map[T]bool = make(map[T]bool)

	for _, v := range q.values {
		if _, ok := set[v]; !ok {
			set[v] = true
		}
	}
	for _, v := range in {
		delete(set, v)
	}
	r.values = make([]T, len(set))
	i := 0
	for k, _ := range set {
		r.values[i] = k
		i++
	}
	return
}

// Count returns number of elements.
func (q Queryable) Count() (count int, err error) {
	return len(q.values), q.err
}

// CountBy returns number of elements satisfying the provided predicate
// function.
func (q Queryable) CountBy(f func(T) (bool, error)) (c int, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		if ok {
			c++
		}
	}
	return
}

// Any determines whether the query source contains any elements.
func (q Queryable) Any() (exists bool, err error) {
	return len(q.values) > 0, q.err
}

// AnyWith determines whether the query source contains any elements satisfying
// the provided predicate function.
func (q Queryable) AnyWith(f func(T) (bool, error)) (exists bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		if ok {
			exists = true
			return
		}
	}
	return
}

// All determines whether all elements of the query source satisfy the provided
// predicate function.
func (q Queryable) All(f func(T) (bool, error)) (all bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}

	all = true // if no elements, result is true
	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		all = all && ok
	}
	return
}

// Single returns the only one element of the original sequence satisfies the
// provided predicate function if exists, otherwise returns ErrNotSinggle.
func (q Queryable) Single(f func(T) (bool, error)) (single T, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	for _, v := range q.values {
		ok, e := f(v)
		if e != nil {
			err = e
			return
		}
		if ok {
			if single != nil {
				err = ErrNotSingle
				return
			}
			single = v
		}
	}

	if single == nil {
		err = ErrNotSingle
	}

	return
}

// ElementAt returns the element at the specified index i. If i is a negative
// number ErrNegativeParam, if no element exists at i-th index, ErrNoElement
// is returned.
func (q Queryable) ElementAt(i int) (elem T, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if i < 0 {
		err = ErrNegativeParam
		return
	}
	if len(q.values) < i+1 {
		err = ErrNoElement
	} else {
		elem = q.values[i]
	}
	return
}

// ElementAtOrNil returns the element at the specified index i if exists,
// otherwise returns nil. If i is a negative number, ErrNegativeParam is
// returned.
func (q Queryable) ElementAtOrNil(i int) (elem T, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if i < 0 {
		err = ErrNegativeParam
		return
	}
	if len(q.values) > i {
		elem = q.values[i]
	}
	return
}

// First returns the element at first position of the query source if exists.
// If source is empty, ErrNoElement is returned.
func (q Queryable) First() (elem T, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		err = ErrNoElement
	} else {
		elem = q.values[0]
	}
	return
}

// FirstOrNil returns the element at first position of the query source, if
// exists. Otherwise returns nil.
func (q Queryable) FirstOrNil() (elem T, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) > 0 {
		elem = q.values[0]
	}
	return
}

func (q Queryable) firstBy(f func(T) (bool, error)) (elem T, found bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	for _, i := range q.values {
		ok, e := f(i)
		if e != nil {
			err = e
			return
		}
		if ok {
			elem = i
			found = true
			break
		}
	}
	return
}

// FirstBy returns the first element in the query source that satisfies the
// provided predicate. If source is empty, ErrNoElement is returned.
func (q Queryable) FirstBy(f func(T) (bool, error)) (elem T, err error) {
	var found bool
	elem, found, err = q.firstBy(f)

	if err == nil && !found {
		err = ErrNoElement
	}
	return
}

// FirstOrNilBy returns the first element in the query source that satisfies
// the provided predicate, if exists, otherwise nil.
func (q Queryable) FirstOrNilBy(f func(T) (bool, error)) (elem T, err error) {
	elem, found, err := q.firstBy(f)
	if !found {
		elem = nil
	}
	return
}

// Last returns the element at last position of the query source if exists.
// If source is empty, ErrNoElement is returned.
func (q Queryable) Last() (elem T, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		err = ErrNoElement
	} else {
		elem = q.values[len(q.values)-1]
	}
	return
}

// LastOrNil returns the element at last index of the query source, if exists.
// Otherwise returns nil.
func (q Queryable) LastOrNil() (elem T, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) > 0 {
		elem = q.values[len(q.values)-1]
	}
	return
}

func (q Queryable) lastBy(f func(T) (bool, error)) (elem T, found bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	for i := len(q.values) - 1; i >= 0; i-- {
		item := q.values[i]
		ok, e := f(item)
		if e != nil {
			err = e
			return
		}
		if ok {
			elem = item
			found = true
			break
		}
	}
	return
}

// LastBy returns the last element in the query source that satisfies the
// provided predicate. If source is empty, ErrNoElement is returned.
func (q Queryable) LastBy(f func(T) (bool, error)) (elem T, err error) {
	var found bool
	elem, found, err = q.lastBy(f)

	if err == nil && !found {
		err = ErrNoElement
	}
	return
}

// LastOrNilBy returns the last element in the query source that satisfies
// the provided predicate, if exists, otherwise nil.
func (q Queryable) LastOrNilBy(f func(T) (bool, error)) (elem T, err error) {
	elem, found, err := q.lastBy(f)
	if !found {
		elem = nil
	}
	return
}

// Reverse returns a query with a inverted order of the original source
func (q Queryable) Reverse() (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	c := len(q.values)
	j := 0
	r.values = make([]T, c)
	for i := c - 1; i >= 0; i-- {
		r.values[j] = q.values[i]
		j++
	}
	return
}

// Take returns a new query with n first elements are taken from the original
// sequence.
func (q Queryable) Take(n int) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if n < 0 {
		n = 0
	}
	if n >= len(q.values) {
		n = len(q.values)
	}
	r.values = q.values[:n]
	return
}

// TakeWhile returns a new query with elements from the original sequence
// by testing them with provided predicate f and stops taking them first time
// predicate returns false.
func (q Queryable) TakeWhile(f func(T) (bool, error)) (r Queryable) {
	n, err := q.findWhileTerminationIndex(f)
	if err != nil {
		r.err = err
		return
	}
	return q.Take(n)
}

// Skip returns a new query with nbypassed
// from the original sequence and takes rest of the elements.
func (q Queryable) Skip(n int) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if n < 0 {
		n = 0
	}
	if n >= len(q.values) {
		n = len(q.values)
	}
	r.values = q.values[n:]
	return
}

// SkipWhile returns a new query with original sequence bypassed
// as long as a provided predicate is true and then takes the
// remaining elements.
func (q Queryable) SkipWhile(f func(T) (bool, error)) (r Queryable) {
	n, err := q.findWhileTerminationIndex(f)
	if err != nil {
		r.err = err
		return
	}
	return q.Skip(n)
}

func (q Queryable) findWhileTerminationIndex(f func(T) (bool, error)) (n int, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if f == nil {
		err = ErrNilFunc
		return
	}
	n = 0
	for _, v := range q.values {
		ok, e := f(v)
		if e != nil {
			err = e
			return
		}
		if ok {
			n++
		} else {
			break
		}
	}
	return
}

// OrderInts returns a new query by sorting integers in the original
// sequence in ascending order. Elements of the original sequence should only be
// int. Otherwise, ErrTypeMismatch will be returned.
func (q Queryable) OrderInts() (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}

	vals, err := toInts(q.values)
	if err != nil {
		r.err = err
		return
	}
	sort.Ints(vals)
	r.values = intsToInterface(vals)

	return
}

// OrderStrings returns a new query by sorting integers in the original
// sequence in ascending order. Elements of the original sequence should only be
// string. Otherwise, ErrTypeMismatch will be returned.
func (q Queryable) OrderStrings() (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	vals, err := toStrings(q.values)
	if err != nil {
		r.err = err
		return
	}
	sort.Strings(vals)
	r.values = stringsToInterface(vals)
	return
}

// OrderFloat64s returns a new query by sorting integers in the original
// sequence in ascending order. Elements of the original sequence should only be
// float64. Otherwise, ErrTypeMismatch will be returned.
func (q Queryable) OrderFloat64s() (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	vals, err := toFloat64s(q.values)
	if err != nil {
		r.err = err
		return
	}
	sort.Float64s(vals)
	r.values = float64sToInterface(vals)
	return
}

// OrderBy returns a new query by sorting elements with provided less function
// in ascending order.
// The comparer function should return true if the parameter "this" is less
// than "that".
func (q Queryable) OrderBy(less func(this T, that T) bool) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if less == nil {
		r.err = ErrNilFunc
		return
	}

	sortQ := sortableQueryable{}
	sortQ.less = less
	sortQ.values = make([]T, len(q.values))
	_ = copy(sortQ.values, q.values)
	sort.Sort(sortQ)
	r.values = sortQ.values
	return
}

// Join correlates the elements of two sequences based on the equality of keys.
// Inner and outer keys are matched using default equality comparer, ==.
//
// Outer sequence is the original sequence.
// Inner sequence is the one provided as input.
// outerKeySelector extracts a key from outer element for comparison.
// innerKeySelector extracts a key from outer element for comparison.
// resultSelector takes key of inner element and key of outer element as input
// and returns a value and these values are returned as a new query.
func (q Queryable) Join(innerCollection []T,
	outerKeySelector func(T) T,
	innerKeySelector func(T) T,
	resultSelector func(
		outer T,
		inner T) T) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if innerCollection == nil {
		r.err = ErrNilInput
		return
	}
	if outerKeySelector == nil || innerKeySelector == nil || resultSelector == nil {
		r.err = ErrNilFunc
		return
	}
	var outerCollection = q.values
	innerKeyLookup := make(map[T]T)

	for _, outer := range outerCollection {
		outerKey := outerKeySelector(outer)
		for _, inner := range innerCollection {
			innerKey, ok := innerKeyLookup[inner]
			if !ok {
				innerKey = innerKeySelector(inner)
				innerKeyLookup[inner] = innerKey
			}
			if innerKey == outerKey {
				elem := resultSelector(outer, inner)
				r.values = append(r.values, elem)
			}
		}
	}
	return
}

// GroupJoin correlates the elements of two sequences based on equality of keys
// and groups the results. The default equality comparer is used to compare
// keys.
//
// Inner and outer keys are matched using default equality comparer, ==.
// Outer sequence is the original sequence.
// Inner sequence is the one provided as input.
// outerKeySelector extracts a key from outer element for comparison.
// innerKeySelector extracts a key from outer element for comparison.
// resultSelector takes key of inner element and key of outer element as input
// and returns a value and these values are returned as a new query.
func (q Queryable) GroupJoin(innerCollection []T,
	outerKeySelector func(T) T,
	innerKeySelector func(T) T,
	resultSelector func(
		outer T,
		inners []T) T) (r Queryable) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if innerCollection == nil {
		r.err = ErrNilInput
		return
	}
	if outerKeySelector == nil || innerKeySelector == nil || resultSelector == nil {
		r.err = ErrNilFunc
		return
	}
	var outerCollection = q.values
	innerKeyLookup := make(map[T]T)

	var results = make(map[T][]T) // outer --> inner...
	for _, outer := range outerCollection {
		outerKey := outerKeySelector(outer)
		bucket := make([]T, 0)
		results[outer] = bucket
		for _, inner := range innerCollection {
			innerKey, ok := innerKeyLookup[inner]
			if !ok {
				innerKey = innerKeySelector(inner)
				innerKeyLookup[inner] = innerKey
			}
			if innerKey == outerKey {
				results[outer] = append(results[outer], inner)
			}
		}
	}

	r.values = make([]T, len(results))
	i := 0
	for k, v := range results {
		outer := k
		inners := v
		r.values[i] = resultSelector(outer, inners)
		i++
	}
	return
}

// Range returns a query with sequence of integral numbers within
// the specified range. int overflows are not handled.
func Range(start, count int) (q Queryable) {
	if count < 0 {
		q.err = ErrNegativeParam
		return
	}
	q.values = make([]T, count)
	for i := 0; i < count; i++ {
		q.values[i] = start + i
	}
	return
}

// Sum computes sum of numeric values in the original sequence.
// See golang spec for numeric types. If sequence has non-numeric types or nil,
// ErrNan is returned.
//
// This method has a poor performance due to language limitations.
// On every element, type assertion is made to find the correct type of the
// element.
func (q Queryable) Sum() (sum float64, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	sum, err = sum_(q.values)
	return
}

func sum_(in []T) (sum float64, err error) {
	// here we do a poor performance operation
	// we use type assertion to convert every numeric value type
	// into float64 for each element in values list
	for i := 0; i < len(in); i++ {
		v := in[i]
		// current optimizations:
		// 1. start from more commonly used types so it terminates early
		if f, ok := v.(int); ok {
			sum += float64(f)
		} else if f, ok := v.(uint); ok {
			sum += float64(f)
		} else if f, ok := v.(float64); ok {
			sum += float64(f)
		} else if f, ok := v.(int32); ok {
			sum += float64(f)
		} else if f, ok := v.(int64); ok {
			sum += float64(f)
		} else if f, ok := v.(float32); ok {
			sum += float64(f)
		} else if f, ok := v.(int8); ok {
			sum += float64(f)
		} else if f, ok := v.(int16); ok {
			sum += float64(f)
		} else if f, ok := v.(uint64); ok {
			sum += float64(f)
		} else if f, ok := v.(uint32); ok {
			sum += float64(f)
		} else if f, ok := v.(uint16); ok {
			sum += float64(f)
		} else if f, ok := v.(uint8); ok {
			sum += float64(f)
		} else {
			err = ErrNan
			return
		}
	}
	return
}

// Average computes average of numeric values in the original sequence.
// See golang spec for numeric types. If sequence has non-numeric types or nil,
// ErrNan is returned. If original sequence is empty, ErrEmptySequence is
// returned.
//
// This method has a poor performance due to language limitations.
// On every element, type assertion is made to find the correct type of the
// element.
func (q Queryable) Average() (avg float64, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	sum, err := sum_(q.values)
	if err != nil {
		return
	}
	avg = sum / float64(len(q.values))
	return
}

// MinInt returns the element with smallest value in the leftmost of the
// sequence. Elements of the original sequence should only be int or
// ErrTypeMismatch is returned. If the sequence is empty ErrEmptySequence is
// returned.
func (q Queryable) MinInt() (min int, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	minIndex, _, err := minMaxInts(q.values)
	if err != nil {
		return
	}
	return q.values[minIndex].(int), nil
}

// MinUint returns the element with smallest value in the leftmost of the
// sequence. Elements of the original sequence should only be uint or
// ErrTypeMismatch is returned. If the sequence is empty ErrEmptySequence is
// returned.
func (q Queryable) MinUint() (min uint, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	minIndex, _, err := minMaxUints(q.values)
	if err != nil {
		return
	}
	return q.values[minIndex].(uint), nil
}

// MinFloat64 returns the element with smallest value in the leftmost of the
// sequence. Elements of the original sequence should only be float64 or
// ErrTypeMismatch is returned. If the sequence is empty ErrEmptySequence is
// returned.
func (q Queryable) MinFloat64() (min float64, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	minIndex, _, err := minMaxFloat64s(q.values)
	if err != nil {
		return
	}
	return q.values[minIndex].(float64), nil
}

// MaxInt returns the element with biggest value in the leftmost of the
// sequence. Elements of the original sequence should only be int or
// ErrTypeMismatch is returned. If the sequence is empty ErrEmptySequence is
// returned.
func (q Queryable) MaxInt() (min int, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	_, maxIndex, err := minMaxInts(q.values)
	if err != nil {
		return
	}
	return q.values[maxIndex].(int), nil
}

// MaxUint returns the element with biggest value in the leftmost of the
// sequence. Elements of the original sequence should only be uint or
// ErrTypeMismatch is returned. If the sequence is empty ErrEmptySequence is
// returned.
func (q Queryable) MaxUint() (min uint, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	_, maxIndex, err := minMaxUints(q.values)
	if err != nil {
		return
	}
	return q.values[maxIndex].(uint), nil
}

// MaxFloat64 returns the element with biggest value in the leftmost of the
// sequence. Elements of the original sequence should only be float64 or
// ErrTypeMismatch is returned. If the sequence is empty ErrEmptySequence is
// returned.
func (q Queryable) MaxFloat64() (min float64, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	_, maxIndex, err := minMaxFloat64s(q.values)
	if err != nil {
		return
	}
	return q.values[maxIndex].(float64), nil
}
