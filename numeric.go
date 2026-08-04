package linq

import (
	"cmp"
	"math"
)

// Number is a constraint that permits any numeric type.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Sum computes the sum of a collection of numeric values. The result has the
// same type as the collection elements. Sum returns zero if the collection
// contains no elements.
func Sum[TNumber Number](q Query[TNumber]) TNumber {
	var sum TNumber
	for item := range q.Iterate {
		sum += item
	}
	return sum
}

// SumBy computes the sum of the values obtained by invoking the selector
// function on each element of the collection. It returns zero if the
// collection contains no elements.
//
// SumBy is a generic method: the numeric result type TNumber is inferred from the
// selector function.
func (q Query[T]) SumBy[TNumber Number](selector func(T) TNumber) TNumber {
	var sum TNumber
	for item := range q.Iterate {
		sum += selector(item)
	}
	return sum
}

// Average computes the average of a collection of numeric values.
// It returns math.NaN() if the collection is empty.
func Average[TNumber Number](q Query[TNumber]) float64 {
	var sum float64
	n := 0
	for item := range q.Iterate {
		sum += float64(item)
		n++
	}

	if n == 0 {
		return math.NaN()
	}
	return sum / float64(n)
}

// AverageBy computes the average of the values obtained by invoking the
// selector function on each element of the collection.
// It returns math.NaN() if the collection is empty.
//
// AverageBy is a generic method: the numeric type TNumber is inferred from the
// selector function.
func (q Query[T]) AverageBy[TNumber Number](selector func(T) TNumber) float64 {
	var sum float64
	n := 0
	for item := range q.Iterate {
		sum += float64(selector(item))
		n++
	}

	if n == 0 {
		return math.NaN()
	}
	return sum / float64(n)
}

// Max returns the maximum value in a collection of ordered values and a
// boolean reporting whether the collection was non-empty.
func Max[T cmp.Ordered](q Query[T]) (T, bool) {
	var r T
	found := false
	for item := range q.Iterate {
		if !found || cmp.Compare(item, r) > 0 {
			r = item
			found = true
		}
	}
	return r, found
}

// Min returns the minimum value in a collection of ordered values and a
// boolean reporting whether the collection was non-empty.
func Min[T cmp.Ordered](q Query[T]) (T, bool) {
	var r T
	found := false
	for item := range q.Iterate {
		if !found || cmp.Compare(item, r) < 0 {
			r = item
			found = true
		}
	}
	return r, found
}

// MaxBy returns the element of a collection with the maximum key, where the
// key is obtained by invoking the selector function on each element, and a
// boolean reporting whether the collection was non-empty.
//
// MaxBy is a generic method: the key type TKey is inferred from the selector
// function and must be an ordered type (cmp.Ordered).
func (q Query[T]) MaxBy[TKey cmp.Ordered](selector func(T) TKey) (T, bool) {
	var r T
	var rKey TKey
	found := false
	for item := range q.Iterate {
		key := selector(item)
		if !found || cmp.Compare(key, rKey) > 0 {
			r, rKey = item, key
			found = true
		}
	}
	return r, found
}

// MinBy returns the element of a collection with the minimum key, where the
// key is obtained by invoking the selector function on each element, and a
// boolean reporting whether the collection was non-empty.
//
// MinBy is a generic method: the key type TKey is inferred from the selector
// function and must be an ordered type (cmp.Ordered).
func (q Query[T]) MinBy[TKey cmp.Ordered](selector func(T) TKey) (T, bool) {
	var r T
	var rKey TKey
	found := false
	for item := range q.Iterate {
		key := selector(item)
		if !found || cmp.Compare(key, rKey) < 0 {
			r, rKey = item, key
			found = true
		}
	}
	return r, found
}
