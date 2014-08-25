// Package linq provides methods for querying and manipulating slices and
// collections.
//
// Author: Ahmet Alp Balkan
package linq

import (
	"errors"
	"sort"
)

// T is an alias for interface{} to make your LINQ code shorter.
type T interface{}

// Query is the type returned from query functions. To evaluate
// get the results of the query, use Results().
type Query struct {
	values []T
	err    error
}

type queryable interface {
	Results() (T, error)
}

type sortableQuery struct {
	values []T
	less   func(this, that T) bool
}

func (q sortableQuery) Len() int           { return len(q.values) }
func (q sortableQuery) Swap(i, j int)      { q.values[i], q.values[j] = q.values[j], q.values[i] }
func (q sortableQuery) Less(i, j int) bool { return q.less(q.values[i], q.values[j]) }

var (
	// ErrNilFunc returned when a predicate, selector or comparer is nil
	ErrNilFunc = errors.New("linq: passed evaluation function is nil")

	// ErrNilInput returned when nil value passed in a slice parameter
	ErrNilInput = errors.New("linq: nil sequence passed as input to function")

	// ErrInvalidInput returned when a non-slice input passed to functions requiring a slice (e.g From, Union, Intersect, Except, Join, GroupJoin)
	ErrInvalidInput = errors.New("linq: non-slice value passed to a T parameter indicating a slice")

	// ErrEmptySequence requested operation is invalid on empty sequences
	ErrEmptySequence = errors.New("linq: empty sequence, operation requires non-empty results sequence")

	// ErrNegativeParam requested when negative value passed to a non-negative requiring parameter
	ErrNegativeParam = errors.New("linq: parameter cannot be negative")

	// ErrNan requested when sequence has invalid elements that method cannot assert into one of builtin numeric types
	ErrNan = errors.New("linq: sequence contains an element of non-numeric types")

	// ErrTypeMismatch returned when sequence nils or elements of different types than function can work with
	ErrTypeMismatch = errors.New("linq: sequence contains element(s) with type different than requested type or nil")

	// ErrNotSingle returned when sequence contains more than one element satisfying the predicate function
	ErrNotSingle = errors.New("linq: sequence contains more than one element matching the given predicate found")
)

// From initializes a linq query with passed slice as the source.
// input parameter must be a slice. Otherwise returns ErrInvalidInput in
// query result evaluating method.
//
// Example:
//     i := []int{1,2,3,4,5,6}
//     q := From(i)
//
// Note: Although it looks like a T (interface{}) input, you should pass a
// slice of any type. There is a hack there to accept any type of slice, which
// is a workaround of type system of Go.
func From(input T) Query {
	var e error
	if input == nil {
		e = ErrNilInput
	}

	out, ok := takeSliceArg(input)
	if !ok {
		e = ErrInvalidInput
		out = nil
	}

	return Query{values: out, err: e}
}

// Results evaluates the query and returns the results as T slice.
// An error occurred in during evaluation of the query will be returned.
//
// Example:
// 	results, err := From(slice).Select(something).Results()
func (q Query) Results() ([]T, error) {
	// we need to copy results for isolating user modification on returned
	// slice from the current Query instance.
	res := make([]T, len(q.values))
	_ = copy(res, q.values)
	return res, q.err
}

// AsParallel returns a ParallelQuery from the same source where the query functions
//  can be executed in parallel for each element of the source with goroutines.
//
// This is an abstraction to not to break user code. If the query method you are
// looking for is not available on ParallelQuery, you can go back to serialized
// Query using AsSequential() method.
func (q Query) AsParallel() ParallelQuery {
	return ParallelQuery{values: q.values, err: q.err}
}

// Where filters a sequence of values based on a predicate function. This
// function will take elements of the source (or results of previous query)
// as interface[] so it should make type assertion to work on the types.
// Returns a query with elements satisfy the condition.
//
// Example:
// 	flying, err := From(animals).Where(func (a T) (bool, error){
//		return a.(*Animal).IsFlying, nil
// 	}).Results()
func (q Query) Where(f func(T) (bool, error)) (r Query) {
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
	return
}

// Select projects each element of a sequence into a new form.
// Returns a query with the result of invoking the transform function
// on each element of original source.
//
// Example:
// 	names, err := From(animals).Select(func (a T) (T, error){
//		return a.(*Animal).Name, nil
// 	}).Results()
func (q Query) Select(f func(T) (T, error)) (r Query) {
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
//
// Example:
// 	distinctInts, err := From(integers).Distinct().Results()
func (q Query) Distinct() (r Query) {
	return q.distinct(nil)
}

// DistinctBy returns distinct elements from the provided source using the
// provided equality comparer. This is a set operation and returns an unordered
// sequence. Number of calls to f will be at most N^2 (all elements are
// distinct) and at best N (all elements are the same).
//
// Example:
// 	distinctFirstNames, err := From(people).DistinctBy(func (p1, p2 T) (bool, error){
//		return  p.(*Person).FirstName == p2.(*Person).FirstName, nil
// 	}).Results()
func (q Query) DistinctBy(f func(T, T) (bool, error)) (r Query) {
	if f == nil {
		r.err = ErrNilFunc
		return
	}
	return q.distinct(f)
}

// distinct returns distinct elements from the provided source using default
// equality comparer (==) or a custom equality comparer function. Complexity
// is O(N).
func (q Query) distinct(f func(T, T) (bool, error)) (r Query) {
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
		for key := range dict {
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
// returns an unordered sequence. inputSlice must be slice of a type although
// it looks like T, otherwise returns ErrInvalidInput.
//
// Example:
// 	all, err := From([]int{1,2,3,4,5}).Union([]int{3,4,5,6}).Results()
// 	// all is {1,2,3,4,5,6}
func (q Query) Union(inputSlice T) (r Query) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if inputSlice == nil {
		r.err = ErrNilInput
		return
	}

	in, ok := takeSliceArg(inputSlice)
	if !ok {
		r.err = ErrInvalidInput
		return
	}

	set := make(map[T]bool)
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
	for k := range set {
		r.values[i] = k
		i++
	}
	return
}

// Intersect returns set intersection of the source sequence and the
// provided input slice using default equality comparer. This is a set
// operation and may return an unordered sequence. inputSlice must be slice of
// a type although it looks like T, otherwise returns ErrInvalidInput.
//
// Example:
// 	both, err := From([]int{1,2,3,4,5}).Intersect([]int{3,4,5,6}).Results()
// 	// both is {3,4,5}
func (q Query) Intersect(inputSlice T) (r Query) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if inputSlice == nil {
		r.err = ErrNilInput
		return
	}

	in, ok := takeSliceArg(inputSlice)
	if !ok {
		r.err = ErrInvalidInput
		return
	}

	set := make(map[T]bool)
	intersection := make(map[T]bool)

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
	for k := range intersection {
		r.values[i] = k
		i++
	}
	return
}

// Except returns set difference of the source sequence and the
// provided input slice using default equality comparer. This is a set
// operation and returns an unordered sequence. inputSlice must be slice of
// a type although it looks like T, otherwise returns ErrInvalidInput.
//
// Example:
// Example:
// 	diffAB, err := From([]int{1,2,3,4,5}).Except([]int{3,4,5,6}).Results()
// 	// diffAB is {1,2}
func (q Query) Except(inputSlice T) (r Query) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if inputSlice == nil {
		r.err = ErrNilInput
		return
	}

	in, ok := takeSliceArg(inputSlice)
	if !ok {
		r.err = ErrInvalidInput
		return
	}

	set := make(map[T]bool)
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
	for k := range set {
		r.values[i] = k
		i++
	}
	return
}

// Count returns number of elements in the sequence. Range().
//
// Example:
// 	over18, err := From(students).Where(func (s T)(bool, error){
//		return s.(*Person).Age > 18, nil
//	}).Count()
func (q Query) Count() (count int, err error) {
	return len(q.values), q.err
}

// CountBy returns number of elements satisfying the provided predicate
// function.
//
// Example:
// 	over18, err := From(students).CountBy(func (s T)(bool, error){
//		return s.(*Person).Age > 18, nil
//	})
func (q Query) CountBy(f func(T) (bool, error)) (c int, err error) {
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
//
// Example:
// 	anyOver18, err := From(students).Where(func (s T)(bool, error){
//		return s.(*Person).Age > 18, nil
//	}).Any()
func (q Query) Any() (exists bool, err error) {
	return len(q.values) > 0, q.err
}

// AnyWith determines whether the query source contains any elements satisfying
// the provided predicate function.
//
// Example:
// 	anyOver18, err := From(students).AnyWith(func (s T)(bool, error){
//		return s.(*Person).Age > 18, nil
//	})
func (q Query) AnyWith(f func(T) (bool, error)) (exists bool, err error) {
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
//
// Returns early if one element does not meet the conditions provided.
//
// Example:
// 	allOver18, err := From(students).All(func (s T)(bool, error){
//		return s.(*Person).Age > 18, nil
//	})
func (q Query) All(f func(T) (bool, error)) (all bool, err error) {
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
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// Single returns the only one element of the original sequence satisfies the
// provided predicate function if exists, otherwise returns ErrNotSingle.
//
// Example:
// 	admin, err := From(students).Single(func (s T)(bool, error){
//		return s.(*Person).Id == 1, nil
//	})
//	if err == nil {
//		// use admin.(*Person)
// 	}
func (q Query) Single(f func(T) (bool, error)) (single T, err error) {
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
// number ErrNegativeParam, if no element exists at i-th index, found will
// be returned false.
//
// Example:
// 	second, found, err := From(ints).OrderInts().ElementAt(2)
// 	if err == nil && found {
//		// use second.(int)
// 	}
func (q Query) ElementAt(i int) (elem T, found bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if i < 0 {
		err = ErrNegativeParam
		return
	}
	if len(q.values) >= i+1 {
		elem = q.values[i]
		found = true
	}
	return
}

// First returns the element at first position of the query source if exists.
// If source is empty or such element is not found, found value will be false,
// otherwise elem is provided.
//
// Example:
// 	first, found, err := From(ints).OrderInts().First()
// 	if err == nil && found {
//		// use first.(int)
// 	}
func (q Query) First() (elem T, found bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) > 0 {
		found = true
		elem = q.values[0]
	}
	return
}

// FirstBy returns the first element in the query source that satisfies the
// provided predicate. If source is empty or such element is not found, found
// value will be false, otherwise elem is provided.
// Example:
// 	second, found, err := From(ints).OrderInts().FirstBy(func (i T)(bool, error){
//		return i.(int) % 2 == 0, nil
// 	})
// 	if err == nil && found {
//		// use first.(int)
// 	}
func (q Query) FirstBy(f func(T) (bool, error)) (elem T, found bool, err error) {
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

// Last returns the element at last position of the query source if exists.
// If source is empty or such element is not found, found
// value will be false, otherwise elem is provided.
//
// Example:
// 	last, found, err := From(ints).OrderInts().Last()
// 	if err == nil && found {
//		// use last.(int)
// 	}
func (q Query) Last() (elem T, found bool, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) > 0 {
		found = true
		elem = q.values[len(q.values)-1]
	}
	return
}

// LastBy returns the last element in the query source that satisfies the
// provided predicate. If source is empty or such element is not found, found
// value will be false, otherwise elem is provided.
//
// Example:
// 	last, found, err := From(ints).OrderInts().LastBy(func (i T)(bool, error){
//		return i.(int) % 2 == 0, nil
// 	})
// 	if err == nil && found {
//		// use last.(int)
// 	}
func (q Query) LastBy(f func(T) (bool, error)) (elem T, found bool, err error) {
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

// Reverse returns a query with a inverted order of the original source
//
// Example:
// 	reversed, err := From([]int{1,2,3,4,5}).Reverse().Results()
func (q Query) Reverse() (r Query) {
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
//
// Example:
// 	arr, err := From([]int{1,2,3,4,5}).Take(3).Results()
//	if err == nil {
//		// arr will be 1, 2, 3
// 	}
func (q Query) Take(n int) (r Query) {
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
// by testing them with provided predicate f and stops taking them first
// predicate returns false.
//
// Example:
// 	arr, err := From([]int{40,10,50,60,100})
// 		.OrderInts()
//		.TakeWhile(func (i T)(bool, error){
//			return i.(int) <= 50, nil
// 		}).Results()
//	if err == nil {
//		// arr will be 10, 40, 50
// 	}
func (q Query) TakeWhile(f func(T) (bool, error)) (r Query) {
	n, err := q.findWhileTerminationIndex(f)
	if err != nil {
		r.err = err
		return
	}
	return q.Take(n)
}

// Skip returns a new query with nbypassed
// from the original sequence and takes rest of the elements.
//
// Example:
// 	arr, err := From([]int{1,2,3,4,5}).Skip(3).Results()
//	if err == nil {
//		// arr will be 4, 5
// 	}
func (q Query) Skip(n int) (r Query) {
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
//
// Example:
// 	arr, err := From([]int{40,10,50,60,100})
// 		.OrderInts()
//		.SkipWhile(func (i T)(bool, error){
//			return i.(int) < 50, nil
// 		}).Results()
//	if err == nil {
//		// arr will be 50, 60, 100
// 	}
func (q Query) SkipWhile(f func(T) (bool, error)) (r Query) {
	n, err := q.findWhileTerminationIndex(f)
	if err != nil {
		r.err = err
		return
	}
	return q.Skip(n)
}

func (q Query) findWhileTerminationIndex(f func(T) (bool, error)) (n int, err error) {
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
//
// Example:
//	sorted, err := From([]int{6,1,2,-1,10}).OrderInts().Results()
//	// sorted = {-1, 1, 2, 6, 10}
func (q Query) OrderInts() (r Query) {
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
//
// Example:
//	sorted, err := From([]string{"foo", "bar", "", "baz"}).OrderStrings().Results()
//	// sorted = {"", "bar", "baz", "foo"}
func (q Query) OrderStrings() (r Query) {
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
//
// Example:
//	sorted, err := From([]float64{-1e-9, -1, 1e-9, 1}}).OrderFloat64s().Results()
//	// sorted = {-1, -1e-9, 1e-9, 1}
func (q Query) OrderFloat64s() (r Query) {
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
//
// Example:
//	sorted, err := From(people).OrderBy(func (this T, that T) bool {
//		return this.(*Person).Age < that.(*Person).Age
// 	}).Results()
func (q Query) OrderBy(less func(this T, that T) bool) (r Query) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if less == nil {
		r.err = ErrNilFunc
		return
	}

	sortQ := sortableQuery{}
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
// Outer collection is the original sequence.
//
// Inner collection is the one provided as innerSlice input parameter as slice
// of any type. Otherwise returns ErrInvalidInput.
//
// outerKeySelector extracts a key from outer element for comparison.
//
// innerKeySelector extracts a key from outer element for comparison.
//
// resultSelector takes outer element and inner element as inputs
// and returns a value which will be an element in the resulting query.
func (q Query) Join(innerSlice T,
	outerKeySelector func(T) T,
	innerKeySelector func(T) T,
	resultSelector func(
		outer T,
		inner T) T) (r Query) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if innerSlice == nil {
		r.err = ErrNilInput
		return
	}
	innerCollection, ok := takeSliceArg(innerSlice)
	if !ok {
		r.err = ErrInvalidInput
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
//
// Outer collection is the original sequence.
//
// Inner collection is the one provided as innerSlice input parameter as slice
// of any type. Otherwise returns ErrInvalidInput.
//
// outerKeySelector extracts a key from outer element for comparison.
//
// innerKeySelector extracts a key from outer element for comparison.
//
// resultSelector takes outer element and inner element as inputs
// and returns a value which will be an element in the resulting query.
func (q Query) GroupJoin(innerSlice T,
	outerKeySelector func(T) T,
	innerKeySelector func(T) T,
	resultSelector func(
		outer T,
		inners []T) T) (r Query) {
	if q.err != nil {
		r.err = q.err
		return
	}
	if innerSlice == nil {
		r.err = ErrNilInput
		return
	}
	innerCollection, ok := takeSliceArg(innerSlice)
	if !ok {
		r.err = ErrInvalidInput
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

// GroupBy returns a map with your specific key and value.
//
// keySelector extracts a key from slice for map key
//
// valueSelector extracts a key from slice for map value
func (q Query) GroupBy(keySelector func(T) T, valueSelector func(T) T) (map[T][]T, error) {
	var results = make(map[T][]T)
	var collection = q.values
	for _, c := range collection {
		key := keySelector(c)
		results[key] = append(results[key], valueSelector(c))
	}
	return results, nil
}

// Range returns a query with sequence of integral numbers within
// the specified range. int overflows are not handled.
//
// Example:
//	seq, err := From(people).Range(-2, 5).Results()
// 	// seq is {-2, -1, 0, 1, 2}
func Range(start, count int) (q Query) {
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
//
// Example:
//	mixedArr = {}interface[]{1, int8(2), uint(3), float64(4.4)}
//	sum, err := From(mixedArr).Sum() // sum is 10.4
//	// or
//	sum, err := From([]int{1,2,3,4,5}).Sum() // sum is 15.0
func (q Query) Sum() (sum float64, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	sum, err = sumMixed(q.values)
	return
}

func sumMixed(in []T) (sum float64, err error) {
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
//
// Example:
//	mixedArr = {}interface[]{1, int8(2), uint(3), float64(4.4)}
//	avg, err := From(mixedArr).Average() // avg is 2.6
//	// or
//	avg, err := From([]int{1,2,3,4,5}).Average() // avg is 3.0
func (q Query) Average() (avg float64, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if len(q.values) == 0 {
		return 0, ErrEmptySequence
	}
	sum, err := sumMixed(q.values)
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
//
// Example:
//	min, err := From([]int{1, -100, 10, 0}).MinInt() // min is -1
func (q Query) MinInt() (min int, err error) {
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
//
// Example:
//	min, err := From([]uint{1, -100, 10, 0}).MinUint() // min is -1
func (q Query) MinUint() (min uint, err error) {
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
//
// Example;
//	min, err := From([]float64{1e-9, 2e10, -1e-10, -1}).MinFloat64() // min is -1.0
func (q Query) MinFloat64() (min float64, err error) {
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
//
// Example:
//	max, err := From([]int{1, -100, 10, 0}).MaxInt() // max is 10
func (q Query) MaxInt() (min int, err error) {
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
//
// Example:
//	max, err := From([]uint{1, -100, 10, 0}).MaxUint() // max is 10
func (q Query) MaxUint() (min uint, err error) {
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
//
// Example:
//	max, err := From([]float64{1e-9, 2e10, -1e-10, -1}).MaxFloat64() // max is 2e10
func (q Query) MaxFloat64() (min float64, err error) {
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
