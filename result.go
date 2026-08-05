package linq

import (
	"iter"
)

// All determines whether all elements of a collection satisfy a condition.
func (q Query[T]) All(predicate func(T) bool) bool {
	for item := range q.Iterate {
		if !predicate(item) {
			return false
		}
	}

	return true
}

// Any determines whether any element of a collection exists.
func (q Query[T]) Any() bool {
	for range q.Iterate {
		return true
	}

	return false
}

// AnyWith determines whether any element of a collection satisfies a condition.
func (q Query[T]) AnyWith(predicate func(T) bool) bool {
	for item := range q.Iterate {
		if predicate(item) {
			return true
		}
	}

	return false
}

// Contains determines whether a collection contains a specified element.
//
// Elements of basic comparable kinds (integers, floats, complex numbers,
// strings, booleans) are compared directly. All other element types are
// compared as boxed (interface) values, so for them this method panics if T
// is not a comparable type at runtime; AnyWith with an equality predicate is
// the fast path for such types.
func (q Query[T]) Contains(value T) bool {
	equal := equalFor[T]()
	for item := range q.Iterate {
		if equal(item, value) {
			return true
		}
	}
	return false
}

// Count returns the number of elements in a collection.
func (q Query[T]) Count() int {
	count := 0
	for range q.Iterate {
		count++
	}
	return count
}

// CountWith returns a number that represents how many elements in the specified
// collection satisfy a condition.
func (q Query[T]) CountWith(predicate func(T) bool) int {
	count := 0
	for item := range q.Iterate {
		if predicate(item) {
			count++
		}
	}
	return count
}

// First returns the first element of a collection and a boolean reporting
// whether the collection was non-empty.
func (q Query[T]) First() (T, bool) {
	for item := range q.Iterate {
		return item, true
	}

	var zero T
	return zero, false
}

// FirstWith returns the first element of a collection that satisfies a
// specified condition and a boolean reporting whether such an element was
// found.
func (q Query[T]) FirstWith(predicate func(T) bool) (T, bool) {
	for item := range q.Iterate {
		if predicate(item) {
			return item, true
		}
	}

	var zero T
	return zero, false
}

// ForEach performs the specified action on each element of a collection.
func (q Query[T]) ForEach(action func(T)) {
	for item := range q.Iterate {
		action(item)
	}
}

// ForEachIndexed performs the specified action on each element of a collection.
//
// The first argument to action represents the zero-based index of that
// element in the source collection. This can be useful if the elements are in a
// known order and you want to do something with an element at a particular
// index, for example. It can also be useful if you want to retrieve the index
// of one or more elements. The second argument to action represents the
// element to process.
func (q Query[T]) ForEachIndexed(action func(int, T)) {
	index := 0
	for item := range q.Iterate {
		action(index, item)
		index++
	}
}

// Last returns the last element of a collection and a boolean reporting
// whether the collection was non-empty.
func (q Query[T]) Last() (T, bool) {
	var r T
	found := false
	for item := range q.Iterate {
		r = item
		found = true
	}

	return r, found
}

// LastWith returns the last element of a collection that satisfies a specified
// condition and a boolean reporting whether such an element was found.
func (q Query[T]) LastWith(predicate func(T) bool) (T, bool) {
	var r T
	found := false
	for item := range q.Iterate {
		if predicate(item) {
			r = item
			found = true
		}
	}

	return r, found
}

// Results collects all items from a query into a slice. It is equivalent to
// ToSlice and is kept for familiarity with earlier go-linq versions.
func (q Query[T]) Results() []T {
	return q.collect()
}

// SequenceEqual determines whether two collections are equal.
//
// Elements of basic comparable kinds (integers, floats, complex numbers,
// strings, booleans) are compared directly. All other element types are
// compared as boxed (interface) values, so for them this method panics if T
// is not a comparable type at runtime.
func (q Query[T]) SequenceEqual(q2 Query[T]) bool {
	next2, stop2 := iter.Pull(q2.Iterate)
	defer stop2()

	eq := equalFor[T]()
	equal := true
	q.Iterate(func(item T) bool {
		item2, ok2 := next2()
		if !ok2 || !eq(item, item2) {
			equal = false
			return false
		}
		return true
	})

	if !equal {
		return false
	}

	_, ok2 := next2()
	return !ok2
}

// Single returns the only element of a collection and a boolean that is true
// only if the collection contains exactly one element.
func (q Query[T]) Single() (T, bool) {
	var r T
	var zero T
	visited := false
	ok := true

	q.Iterate(func(item T) bool {
		if visited {
			ok = false
			return false
		}

		r = item
		visited = true
		return true
	})

	if !visited || !ok {
		return zero, false
	}
	return r, true
}

// SingleWith returns the only element of a collection that satisfies a
// specified condition and a boolean that is true only if exactly one such
// element exists.
func (q Query[T]) SingleWith(predicate func(T) bool) (T, bool) {
	var r T
	var zero T
	found := false
	ok := true

	q.Iterate(func(item T) bool {
		if !predicate(item) {
			return true
		}

		if found {
			ok = false
			return false
		}

		r = item
		found = true
		return true
	})

	if !found || !ok {
		return zero, false
	}
	return r, true
}

// ToChannel iterates over a collection and outputs each element to a channel,
// then closes it.
func (q Query[T]) ToChannel(result chan<- T) {
	defer close(result)

	for item := range q.Iterate {
		result <- item
	}
}

// ToMap iterates over a collection of KeyValue elements and returns a map
// populated with them. To populate a map from a collection of other types,
// use the ToMapBy method.
func ToMap[TKey comparable, TValue any](q Query[KeyValue[TKey, TValue]]) map[TKey]TValue {
	result := make(map[TKey]TValue)
	for item := range q.Iterate {
		result[item.Key] = item.Value
	}
	return result
}

// ToMapBy iterates over a collection and returns a map populated with
// elements. Functions keySelector and valueSelector are executed for each
// element of the collection to generate the key and value for the map.
//
// ToMapBy is a generic method: the map key type TKey and value type TValue are
// inferred from the selector functions. The key type TKey must be comparable.
func (q Query[T]) ToMapBy[TKey comparable, TValue any](
	keySelector func(T) TKey,
	valueSelector func(T) TValue) map[TKey]TValue {
	result := make(map[TKey]TValue)
	for item := range q.Iterate {
		result[keySelector(item)] = valueSelector(item)
	}
	return result
}

// ToSlice iterates over a collection and returns the results as a slice.
// When the query's element count is known up front (e.g. a slice source
// transformed only by length-preserving operators), the result slice is
// allocated once at the right size.
func (q Query[T]) ToSlice() []T {
	return q.collect()
}
