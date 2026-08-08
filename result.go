package linq

import (
	"iter"
)

// All determines whether all elements of a collection satisfy a condition.
func (q Query[T]) All(predicate func(T) bool) bool {
	all := true
	q.Iterate(func(item T) bool {
		if !predicate(item) {
			all = false
			return false
		}
		return true
	})
	return all
}

// Any determines whether any element of a collection exists.
func (q Query[T]) Any() bool {
	any := false
	q.Iterate(func(T) bool {
		any = true
		return false
	})
	return any
}

// AnyWith determines whether any element of a collection satisfies a condition.
func (q Query[T]) AnyWith(predicate func(T) bool) bool {
	found := false
	q.Iterate(func(item T) bool {
		if predicate(item) {
			found = true
			return false
		}
		return true
	})
	return found
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
	found := false
	q.Iterate(func(item T) bool {
		if equal(item, value) {
			found = true
			return false
		}
		return true
	})
	return found
}

// Count returns the number of elements in a collection.
func (q Query[T]) Count() int {
	count := 0
	q.Iterate(func(T) bool {
		count++
		return true
	})
	return count
}

// CountWith returns a number that represents how many elements in the specified
// collection satisfy a condition.
func (q Query[T]) CountWith(predicate func(T) bool) int {
	count := 0
	q.Iterate(func(item T) bool {
		if predicate(item) {
			count++
		}
		return true
	})
	return count
}

// First returns the first element of a collection and a boolean reporting
// whether the collection was non-empty.
func (q Query[T]) First() (T, bool) {
	var r T
	found := false
	q.Iterate(func(item T) bool {
		r = item
		found = true
		return false
	})
	return r, found
}

// FirstWith returns the first element of a collection that satisfies a
// specified condition and a boolean reporting whether such an element was
// found.
func (q Query[T]) FirstWith(predicate func(T) bool) (T, bool) {
	var r T
	found := false
	q.Iterate(func(item T) bool {
		if predicate(item) {
			r = item
			found = true
			return false
		}
		return true
	})
	return r, found
}

// ForEach performs the specified action on each element of a collection.
func (q Query[T]) ForEach(action func(T)) {
	q.Iterate(func(item T) bool {
		action(item)
		return true
	})
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
	q.Iterate(func(item T) bool {
		action(index, item)
		index++
		return true
	})
}

// Last returns the last element of a collection and a boolean reporting
// whether the collection was non-empty.
func (q Query[T]) Last() (T, bool) {
	var r T
	found := false
	q.Iterate(func(item T) bool {
		r = item
		found = true
		return true
	})
	return r, found
}

// LastWith returns the last element of a collection that satisfies a specified
// condition and a boolean reporting whether such an element was found.
func (q Query[T]) LastWith(predicate func(T) bool) (T, bool) {
	var r T
	found := false
	q.Iterate(func(item T) bool {
		if predicate(item) {
			r = item
			found = true
		}
		return true
	})
	return r, found
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

	q.Iterate(func(item T) bool {
		result <- item
		return true
	})
}

// ToMap iterates over a collection of KeyValue elements and returns a map
// populated with them. To populate a map from a collection of other types,
// use the ToMapBy method.
func ToMap[TKey comparable, TValue any](q Query[KeyValue[TKey, TValue]]) map[TKey]TValue {
	result := make(map[TKey]TValue)
	q.Iterate(func(item KeyValue[TKey, TValue]) bool {
		result[item.Key] = item.Value
		return true
	})
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
	q.Iterate(func(item T) bool {
		result[keySelector(item)] = valueSelector(item)
		return true
	})
	return result
}

// ToSlice iterates over a collection and returns the results as a slice.
// When the query's element count is known up front (e.g. a slice source
// transformed only by length-preserving operators), the result slice is
// allocated once at the right size.
func (q Query[T]) ToSlice() []T {
	return q.collect()
}
