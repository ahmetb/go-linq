package linq

import (
	"iter"
	"math"
	"reflect"
	"slices"
)

// All determines whether all elements of a collection satisfy a condition.
func (q Query) All(predicate func(any) bool) bool {
	for item := range q.Iterate {
		if !predicate(item) {
			return false
		}
	}

	return true
}

// AllT is the typed version of All.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: All has better performance than AllT.
func (q Query) AllT(predicateFn any) bool {

	predicateGenericFunc, err := newGenericFunc(
		"AllT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}
	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.All(predicateFunc)
}

// Any determines whether any element of a collection exists.
func (q Query) Any() bool {
	for range q.Iterate {
		return true
	}

	return false
}

// AnyWith determines whether any element of a collection satisfies a condition.
func (q Query) AnyWith(predicate func(any) bool) bool {
	for item := range q.Iterate {
		if predicate(item) {
			return true
		}
	}

	return false
}

// AnyWithT is the typed version of AnyWith.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: AnyWith has better performance than AnyWithT.
func (q Query) AnyWithT(predicateFn any) bool {

	predicateGenericFunc, err := newGenericFunc(
		"AnyWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.AnyWith(predicateFunc)
}

// Average computes the average of a collection of numeric values.
// It panics if the sequence contains non-numeric types.
// It returns math.NaN() if the sequence is empty.
func (q Query) Average() (r float64) {
	next, stop := iter.Pull(q.Iterate)
	defer stop()

	item, ok := next()
	if !ok {
		return math.NaN()
	}

	n := 1
	switch item.(type) {
	case int, int8, int16, int32, int64:
		conv := getIntConverter(item)
		sum := conv(item)

		for item, ok = next(); ok; item, ok = next() {
			sum += conv(item)
			n++
		}

		r = float64(sum)
	case uint, uint8, uint16, uint32, uint64:
		conv := getUIntConverter(item)
		sum := conv(item)

		for item, ok = next(); ok; item, ok = next() {
			sum += conv(item)
			n++
		}

		r = float64(sum)
	default:
		conv := getFloatConverter(item)
		r = conv(item)

		for item, ok = next(); ok; item, ok = next() {
			r += conv(item)
			n++
		}
	}

	return r / float64(n)
}

// Contains determines whether a collection contains a specified element.
func (q Query) Contains(value any) bool {
	for item := range q.Iterate {
		if item == value {
			return true
		}
	}
	return false
}

// Count returns the number of elements in a collection.
func (q Query) Count() int {
	count := 0
	for range q.Iterate {
		count++
	}
	return count
}

// CountWith returns a number that represents how many elements in the specified
// collection satisfy a condition.
func (q Query) CountWith(predicate func(any) bool) int {
	count := 0
	for item := range q.Iterate {
		if predicate(item) {
			count++
		}
	}
	return count
}

// CountWithT is the typed version of CountWith.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: CountWith has better performance than CountWithT.
func (q Query) CountWithT(predicateFn any) int {

	predicateGenericFunc, err := newGenericFunc(
		"CountWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.CountWith(predicateFunc)
}

// First returns the first element of a collection.
func (q Query) First() any {
	for item := range q.Iterate {
		return item
	}

	return nil
}

// FirstWith returns the first element of a collection that satisfies a
// specified condition.
func (q Query) FirstWith(predicate func(any) bool) any {
	for item := range q.Iterate {
		if predicate(item) {
			return item
		}
	}

	return nil
}

// FirstWithT is the typed version of FirstWith.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: FirstWith has better performance than FirstWithT.
func (q Query) FirstWithT(predicateFn any) any {

	predicateGenericFunc, err := newGenericFunc(
		"FirstWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.FirstWith(predicateFunc)
}

// ForEach performs the specified action on each element of a collection.
func (q Query) ForEach(action func(any)) {
	for item := range q.Iterate {
		action(item)
	}
}

// ForEachT is the typed version of ForEach.
//
//   - actionFn is of type "func(TSource)"
//
// NOTE: ForEach has better performance than ForEachT.
func (q Query) ForEachT(actionFn any) {
	actionGenericFunc, err := newGenericFunc(
		"ForEachT", "actionFn", actionFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), nil),
	)

	if err != nil {
		panic(err)
	}

	actionFunc := func(item any) {
		actionGenericFunc.Call(item)
	}

	q.ForEach(actionFunc)
}

// ForEachIndexed performs the specified action on each element of a collection.
//
// The first argument to action represents the zero-based index of that
// element in the source collection. This can be useful if the elements are in a
// known order and you want to do something with an element at a particular
// index, for example. It can also be useful if you want to retrieve the index
// of one or more elements. The second argument to action represents the
// element to process.
func (q Query) ForEachIndexed(action func(int, any)) {
	index := 0
	for item := range q.Iterate {
		action(index, item)
		index++
	}
}

// ForEachIndexedT is the typed version of ForEachIndexed.
//
//   - actionFn is of type "func(int, TSource)"
//
// NOTE: ForEachIndexed has better performance than ForEachIndexedT.
func (q Query) ForEachIndexedT(actionFn any) {
	actionGenericFunc, err := newGenericFunc(
		"ForEachIndexedT", "actionFn", actionFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), nil),
	)

	if err != nil {
		panic(err)
	}

	actionFunc := func(index int, item any) {
		actionGenericFunc.Call(index, item)
	}

	q.ForEachIndexed(actionFunc)
}

// Last returns the last element of a collection.
func (q Query) Last() (r any) {
	for r = range q.Iterate {
	}

	return
}

// LastWith returns the last element of a collection that satisfies a specified
// condition.
func (q Query) LastWith(predicate func(any) bool) (r any) {
	for item := range q.Iterate {
		if predicate(item) {
			r = item
		}
	}

	return
}

// LastWithT is the typed version of LastWith.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: LastWith has better performance than LastWithT.
func (q Query) LastWithT(predicateFn any) any {

	predicateGenericFunc, err := newGenericFunc(
		"LastWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.LastWith(predicateFunc)
}

// Max returns the maximum value in a collection of values.
func (q Query) Max() any {
	next, stop := iter.Pull(q.Iterate)
	defer stop()

	r, ok := next()
	if !ok {
		return nil
	}

	compare := getComparer(r)

	for item, ok := next(); ok; item, ok = next() {
		if compare(item, r) > 0 {
			r = item
		}
	}

	return r
}

// Min returns the minimum value in a collection of values.
func (q Query) Min() any {
	next, stop := iter.Pull(q.Iterate)
	defer stop()

	r, ok := next()
	if !ok {
		return nil
	}

	compare := getComparer(r)

	for item, ok := next(); ok; item, ok = next() {
		if compare(item, r) < 0 {
			r = item
		}
	}

	return r
}

// Results collects all items from a query into a slice.
func (q Query) Results() []any {
	return slices.Collect(q.Iterate)
}

// SequenceEqual determines whether two collections are equal.
func (q Query) SequenceEqual(q2 Query) bool {
	next, stop := iter.Pull(q.Iterate)
	defer stop()

	next2, stop2 := iter.Pull(q2.Iterate)
	defer stop2()

	for item, ok := next(); ok; item, ok = next() {
		item2, ok2 := next2()
		if !ok2 || item != item2 {
			return false
		}
	}

	_, ok2 := next2()
	return !ok2
}

// Single returns the only element of a collection, and nil if there is not
// exactly one element in the collection.
func (q Query) Single() (r any) {
	visited := false
	for item := range q.Iterate {
		if visited {
			return nil
		}

		r = item
		visited = true
	}

	return
}

// SingleWith returns the only element of a collection that satisfies a
// specified condition, and nil if more than one such element exists.
func (q Query) SingleWith(predicate func(any) bool) (r any) {
	found := false
	for item := range q.Iterate {
		if !predicate(item) {
			continue
		}

		if found {
			return nil
		}

		r = item
		found = true

	}

	return
}

// SingleWithT is the typed version of SingleWith.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: SingleWith has better performance than SingleWithT.
func (q Query) SingleWithT(predicateFn any) any {
	predicateGenericFunc, err := newGenericFunc(
		"SingleWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item any) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.SingleWith(predicateFunc)
}

// SumInts computes the sum of a collection of numeric values.
//
// Values can be of any integer type: int, int8, int16, int32, int64. The result
// is int64. Method returns zero if the collection contains no elements.
func (q Query) SumInts() (r int64) {
	next, stop := iter.Pull(q.Iterate)
	defer stop()

	item, ok := next()
	if !ok {
		return 0
	}

	conv := getIntConverter(item)
	r = conv(item)

	for item, ok = next(); ok; item, ok = next() {
		r += conv(item)
	}

	return
}

// SumUInts computes the sum of a collection of numeric values.
//
// Values can be of any unsigned integer type: uint, uint8, uint16, uint32,
// uint64. The result is uint64. Method returns zero if the collection contains no
// elements.
func (q Query) SumUInts() (r uint64) {
	next, stop := iter.Pull(q.Iterate)
	defer stop()

	item, ok := next()
	if !ok {
		return 0
	}

	conv := getUIntConverter(item)
	r = conv(item)

	for item, ok = next(); ok; item, ok = next() {
		r += conv(item)
	}

	return
}

// SumFloats computes the sum of a collection of numeric values.
//
// Values can be of any float type: float32 or float64. The result is float64.
// Method returns zero if the collection contains no elements.
func (q Query) SumFloats() (r float64) {
	next, stop := iter.Pull(q.Iterate)
	defer stop()

	item, ok := next()
	if !ok {
		return 0
	}

	conv := getFloatConverter(item)
	r = conv(item)

	for item, ok = next(); ok; item, ok = next() {
		r += conv(item)
	}

	return
}

// ToChannel iterates over a collection and outputs each element to a channel,
// then closes it.
func (q Query) ToChannel(result chan<- any) {
	defer close(result)

	for item := range q.Iterate {
		result <- item
	}
}

// ToChannelT is the typed version of ToChannel.
//
//   - result is of type "chan TSource"
//
// NOTE: ToChannel has better performance than ToChannelT.
func (q Query) ToChannelT(result any) {
	r := reflect.ValueOf(result)

	for item := range q.Iterate {
		r.Send(reflect.ValueOf(item))
	}

	r.Close()
}

// ToMap iterates over a collection and populates a result map with elements.
// Collection elements have to be of KeyValue type to use this method. To
// populate a map with elements of different types, use the ToMapBy method. ToMap
// doesn't empty the result map before populating it.
func (q Query) ToMap(result any) {
	q.ToMapBy(
		result,
		func(i any) any {
			return i.(KeyValue).Key
		},
		func(i any) any {
			return i.(KeyValue).Value
		})
}

// ToMapBy iterates over a collection and populates the result map with
// elements. Functions keySelector and valueSelector are executed for each
// element of the collection to generate key and value for the map. Generated
// key and value types must be assignable to the map's key and value types.
// ToMapBy doesn't empty the result map before populating it.
func (q Query) ToMapBy(result any,
	keySelector func(any) any,
	valueSelector func(any) any) {
	res := reflect.ValueOf(result)
	m := reflect.Indirect(res)

	for item := range q.Iterate {
		key := reflect.ValueOf(keySelector(item))
		value := reflect.ValueOf(valueSelector(item))

		m.SetMapIndex(key, value)
	}

	res.Elem().Set(m)
}

// ToMapByT is the typed version of ToMapBy.
//
//   - keySelectorFn is of type "func(TSource)TKey"
//   - valueSelectorFn is of type "func(TSource)TValue"
//
// NOTE: ToMapBy has better performance than ToMapByT.
func (q Query) ToMapByT(result any,
	keySelectorFn any, valueSelectorFn any) {
	keySelectorGenericFunc, err := newGenericFunc(
		"ToMapByT", "keySelectorFn", keySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	keySelectorFunc := func(item any) any {
		return keySelectorGenericFunc.Call(item)
	}

	valueSelectorGenericFunc, err := newGenericFunc(
		"ToMapByT", "valueSelectorFn", valueSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	valueSelectorFunc := func(item any) any {
		return valueSelectorGenericFunc.Call(item)
	}

	q.ToMapBy(result, keySelectorFunc, valueSelectorFunc)
}

// ToSlice iterates over a collection and saves the results in the slice pointed
// by v. It overwrites the existing slice, starting from index 0.
//
// If the slice pointed by v has sufficient capacity, v will be pointed to a
// resliced slice. If it does not, a new underlying array will be allocated and
// v will point to it.
func (q Query) ToSlice(v any) {
	ptrValue := reflect.ValueOf(v)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.IsNil() {
		panic("ToSlice: v must be a pointer to a slice")
	}

	sliceValue := reflect.Indirect(ptrValue)
	if sliceValue.Kind() != reflect.Slice {
		panic("ToSlice: v must point to a slice")
	}

	// Reset length to 0 but keep capacity (like s = s[:0])
	// This preserves any existing capacity for reuse.
	out := sliceValue.Slice(0, 0)

	elemType := sliceValue.Type().Elem()
	for item := range q.Iterate {
		itemValue := reflect.ValueOf(item)

		// Ensure type compatibility with the slice element type.
		if !itemValue.Type().AssignableTo(elemType) {
			if itemValue.Type().ConvertibleTo(elemType) {
				itemValue = itemValue.Convert(elemType)
			} else {
				panic("ToSlice: item type is not assignable/convertible to slice element type")
			}
		}

		out = reflect.Append(out, itemValue)
	}

	// Point v to the final slice (which may have a new backing array).
	ptrValue.Elem().Set(out)
}
