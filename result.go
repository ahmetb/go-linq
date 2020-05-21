package linq

import (
	"math"
	"reflect"
)

// All determines whether all elements of a collection satisfy a condition.
func (q Query) All(predicate func(interface{}) bool) bool {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
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
func (q Query) AllT(predicateFn interface{}) bool {

	predicateGenericFunc, err := newGenericFunc(
		"AllT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}
	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.All(predicateFunc)
}

// Any determines whether any element of a collection exists.
func (q Query) Any() bool {
	_, ok := q.Iterate()()
	return ok
}

// AnyWith determines whether any element of a collection satisfies a condition.
func (q Query) AnyWith(predicate func(interface{}) bool) bool {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
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
func (q Query) AnyWithT(predicateFn interface{}) bool {

	predicateGenericFunc, err := newGenericFunc(
		"AnyWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.AnyWith(predicateFunc)
}

// Average computes the average of a collection of numeric values.
func (q Query) Average() (r float64) {
	next := q.Iterate()
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
func (q Query) Contains(value interface{}) bool {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if item == value {
			return true
		}
	}

	return false
}

// Count returns the number of elements in a collection.
func (q Query) Count() (r int) {
	next := q.Iterate()

	for _, ok := next(); ok; _, ok = next() {
		r++
	}

	return
}

// CountWith returns a number that represents how many elements in the specified
// collection satisfy a condition.
func (q Query) CountWith(predicate func(interface{}) bool) (r int) {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if predicate(item) {
			r++
		}
	}

	return
}

// CountWithT is the typed version of CountWith.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: CountWith has better performance than CountWithT.
func (q Query) CountWithT(predicateFn interface{}) int {

	predicateGenericFunc, err := newGenericFunc(
		"CountWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.CountWith(predicateFunc)
}

// First returns the first element of a collection.
func (q Query) First() interface{} {
	item, _ := q.Iterate()()
	return item
}

// FirstWith returns the first element of a collection that satisfies a
// specified condition.
func (q Query) FirstWith(predicate func(interface{}) bool) interface{} {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
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
func (q Query) FirstWithT(predicateFn interface{}) interface{} {

	predicateGenericFunc, err := newGenericFunc(
		"FirstWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.FirstWith(predicateFunc)
}

// ForEach performs the specified action on each element of a collection.
func (q Query) ForEach(action func(interface{})) {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		action(item)
	}
}

// ForEachT is the typed version of ForEach.
//
//   - actionFn is of type "func(TSource)"
//
// NOTE: ForEach has better performance than ForEachT.
func (q Query) ForEachT(actionFn interface{}) {
	actionGenericFunc, err := newGenericFunc(
		"ForEachT", "actionFn", actionFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), nil),
	)

	if err != nil {
		panic(err)
	}

	actionFunc := func(item interface{}) {
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
func (q Query) ForEachIndexed(action func(int, interface{})) {
	next := q.Iterate()
	index := 0

	for item, ok := next(); ok; item, ok = next() {
		action(index, item)
		index++
	}
}

// ForEachIndexedT is the typed version of ForEachIndexed.
//
//   - actionFn is of type "func(int, TSource)"
//
// NOTE: ForEachIndexed has better performance than ForEachIndexedT.
func (q Query) ForEachIndexedT(actionFn interface{}) {
	actionGenericFunc, err := newGenericFunc(
		"ForEachIndexedT", "actionFn", actionFn,
		simpleParamValidator(newElemTypeSlice(new(int), new(genericType)), nil),
	)

	if err != nil {
		panic(err)
	}

	actionFunc := func(index int, item interface{}) {
		actionGenericFunc.Call(index, item)
	}

	q.ForEachIndexed(actionFunc)
}

// Last returns the last element of a collection.
func (q Query) Last() (r interface{}) {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		r = item
	}

	return
}

// LastWith returns the last element of a collection that satisfies a specified
// condition.
func (q Query) LastWith(predicate func(interface{}) bool) (r interface{}) {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
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
func (q Query) LastWithT(predicateFn interface{}) interface{} {

	predicateGenericFunc, err := newGenericFunc(
		"LastWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.LastWith(predicateFunc)
}

// Max returns the maximum value in a collection of values.
func (q Query) Max() (r interface{}) {
	next := q.Iterate()
	item, ok := next()
	if !ok {
		return nil
	}

	compare := getComparer(item)
	r = item

	for item, ok := next(); ok; item, ok = next() {
		if compare(item, r) > 0 {
			r = item
		}
	}

	return
}

// Min returns the minimum value in a collection of values.
func (q Query) Min() (r interface{}) {
	next := q.Iterate()
	item, ok := next()
	if !ok {
		return nil
	}

	compare := getComparer(item)
	r = item

	for item, ok := next(); ok; item, ok = next() {
		if compare(item, r) < 0 {
			r = item
		}
	}

	return
}

// Results iterates over a collection and returnes slice of interfaces
func (q Query) Results() (r []interface{}) {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		r = append(r, item)
	}

	return
}

// SequenceEqual determines whether two collections are equal.
func (q Query) SequenceEqual(q2 Query) bool {
	next := q.Iterate()
	next2 := q2.Iterate()

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
func (q Query) Single() interface{} {
	next := q.Iterate()
	item, ok := next()
	if !ok {
		return nil
	}

	_, ok = next()
	if ok {
		return nil
	}

	return item
}

// SingleWith returns the only element of a collection that satisfies a
// specified condition, and nil if more than one such element exists.
func (q Query) SingleWith(predicate func(interface{}) bool) (r interface{}) {
	next := q.Iterate()
	found := false

	for item, ok := next(); ok; item, ok = next() {
		if predicate(item) {
			if found {
				return nil
			}

			found = true
			r = item
		}
	}

	return
}

// SingleWithT is the typed version of SingleWith.
//
//   - predicateFn is of type "func(TSource) bool"
//
// NOTE: SingleWith has better performance than SingleWithT.
func (q Query) SingleWithT(predicateFn interface{}) interface{} {
	predicateGenericFunc, err := newGenericFunc(
		"SingleWithT", "predicateFn", predicateFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
	)
	if err != nil {
		panic(err)
	}

	predicateFunc := func(item interface{}) bool {
		return predicateGenericFunc.Call(item).(bool)
	}

	return q.SingleWith(predicateFunc)
}

// SumInts computes the sum of a collection of numeric values.
//
// Values can be of any integer type: int, int8, int16, int32, int64. The result
// is int64. Method returns zero if collection contains no elements.
func (q Query) SumInts() (r int64) {
	next := q.Iterate()
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
// uint64. The result is uint64. Method returns zero if collection contains no
// elements.
func (q Query) SumUInts() (r uint64) {
	next := q.Iterate()
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
// Method returns zero if collection contains no elements.
func (q Query) SumFloats() (r float64) {
	next := q.Iterate()
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
func (q Query) ToChannel(result chan<- interface{}) {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		result <- item
	}

	close(result)
}

// ToChannelT is the typed version of ToChannel.
//
//   - result is of type "chan TSource"
//
// NOTE: ToChannel has better performance than ToChannelT.
func (q Query) ToChannelT(result interface{}) {
	r := reflect.ValueOf(result)
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		r.Send(reflect.ValueOf(item))
	}

	r.Close()
}

// ToMap iterates over a collection and populates result map with elements.
// Collection elements have to be of KeyValue type to use this method. To
// populate a map with elements of different type use ToMapBy method. ToMap
// doesn't empty the result map before populating it.
func (q Query) ToMap(result interface{}) {
	q.ToMapBy(
		result,
		func(i interface{}) interface{} {
			return i.(KeyValue).Key
		},
		func(i interface{}) interface{} {
			return i.(KeyValue).Value
		})
}

// ToMapBy iterates over a collection and populates the result map with
// elements. Functions keySelector and valueSelector are executed for each
// element of the collection to generate key and value for the map. Generated
// key and value types must be assignable to the map's key and value types.
// ToMapBy doesn't empty the result map before populating it.
func (q Query) ToMapBy(result interface{},
	keySelector func(interface{}) interface{},
	valueSelector func(interface{}) interface{}) {
	res := reflect.ValueOf(result)
	m := reflect.Indirect(res)
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
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
func (q Query) ToMapByT(result interface{},
	keySelectorFn interface{}, valueSelectorFn interface{}) {
	keySelectorGenericFunc, err := newGenericFunc(
		"ToMapByT", "keySelectorFn", keySelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	keySelectorFunc := func(item interface{}) interface{} {
		return keySelectorGenericFunc.Call(item)
	}

	valueSelectorGenericFunc, err := newGenericFunc(
		"ToMapByT", "valueSelectorFn", valueSelectorFn,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	valueSelectorFunc := func(item interface{}) interface{} {
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
func (q Query) ToSlice(v interface{}) {
	res := reflect.ValueOf(v)
	slice := reflect.Indirect(res)

	cap := slice.Cap()
	res.Elem().Set(slice.Slice(0, cap)) // make len(slice)==cap(slice) from now on

	next := q.Iterate()
	index := 0
	for item, ok := next(); ok; item, ok = next() {
		if index >= cap {
			slice, cap = grow(slice)
		}
		slice.Index(index).Set(reflect.ValueOf(item))
		index++
	}

	// reslice the len(res)==cap(res) actual res size
	res.Elem().Set(slice.Slice(0, index))
}

// grow grows the slice s by doubling its capacity, then it returns the new
// slice (resliced to its full capacity) and the new capacity.
func grow(s reflect.Value) (v reflect.Value, newCap int) {
	cap := s.Cap()
	if cap == 0 {
		cap = 1
	} else {
		cap *= 2
	}
	newSlice := reflect.MakeSlice(s.Type(), cap, cap)
	reflect.Copy(newSlice, s)
	return newSlice, cap
}
