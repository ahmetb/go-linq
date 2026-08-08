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
	q.Iterate(func(item TNumber) bool {
		sum += item
		return true
	})
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
	q.Iterate(func(item T) bool {
		sum += selector(item)
		return true
	})
	return sum
}

// numberKind classifies the members of the Number constraint by how their
// values must be accumulated.
type numberKind int

const (
	signedNumber numberKind = iota
	unsignedNumber
	floatingPointNumber
)

// kindOfNumber reports whether TNumber is a signed integer, unsigned integer,
// or floating-point type. It probes the type's arithmetic behavior rather than
// switching on concrete types so that named types in the constraint's type set
// (e.g. "type ID int") are classified by their underlying kind.
func kindOfNumber[TNumber Number]() numberKind {
	half := 0.5
	if TNumber(half) != 0 {
		return floatingPointNumber // an integer conversion would have truncated to 0
	}
	if TNumber(0)-1 > 0 {
		return unsignedNumber // the subtraction wrapped around to the maximum value
	}
	return signedNumber
}

// The average* helpers below accumulate integer values in 64-bit integer
// arithmetic and convert to float64 only for the final division, so integer
// inputs beyond 2^53 are not rounded before they are summed (matching v4
// semantics). Like v4, integer sums that exceed the int64/uint64 range wrap
// around. Each accumulation loop lives in its own function so that only the
// accumulator of the kind actually iterated escapes to the heap.

func averageSigned[T any, TNumber Number](q Query[T], selector func(T) TNumber) (float64, int) {
	var sum int64
	n := 0
	q.Iterate(func(item T) bool {
		sum += int64(selector(item))
		n++
		return true
	})
	return float64(sum), n
}

func averageUnsigned[T any, TNumber Number](q Query[T], selector func(T) TNumber) (float64, int) {
	var sum uint64
	n := 0
	q.Iterate(func(item T) bool {
		sum += uint64(selector(item))
		n++
		return true
	})
	return float64(sum), n
}

func averageFloat[T any, TNumber Number](q Query[T], selector func(T) TNumber) (float64, int) {
	var sum float64
	n := 0
	q.Iterate(func(item T) bool {
		sum += float64(selector(item))
		n++
		return true
	})
	return sum, n
}

func averageBy[T any, TNumber Number](q Query[T], selector func(T) TNumber) float64 {
	var sum float64
	var n int
	switch kindOfNumber[TNumber]() {
	case signedNumber:
		sum, n = averageSigned(q, selector)
	case unsignedNumber:
		sum, n = averageUnsigned(q, selector)
	default: // floatingPointNumber
		sum, n = averageFloat(q, selector)
	}

	if n == 0 {
		return math.NaN()
	}
	return sum / float64(n)
}

// Average computes the average of a collection of numeric values. Integer
// values are summed in 64-bit integer arithmetic before dividing, so values
// beyond 2^53 do not lose precision to intermediate float64 rounding.
// It returns math.NaN() if the collection is empty.
func Average[TNumber Number](q Query[TNumber]) float64 {
	// The loops are spelled out rather than delegated to averageBy with an
	// identity selector, which would cost an indirect call per element (~2x).
	switch kindOfNumber[TNumber]() {
	case signedNumber:
		var sum int64
		n := 0
		q.Iterate(func(item TNumber) bool {
			sum += int64(item)
			n++
			return true
		})
		if n == 0 {
			return math.NaN()
		}
		return float64(sum) / float64(n)
	case unsignedNumber:
		var sum uint64
		n := 0
		q.Iterate(func(item TNumber) bool {
			sum += uint64(item)
			n++
			return true
		})
		if n == 0 {
			return math.NaN()
		}
		return float64(sum) / float64(n)
	default: // floatingPointNumber
		var sum float64
		n := 0
		q.Iterate(func(item TNumber) bool {
			sum += float64(item)
			n++
			return true
		})
		if n == 0 {
			return math.NaN()
		}
		return sum / float64(n)
	}
}

// AverageBy computes the average of the values obtained by invoking the
// selector function on each element of the collection. Integer values are
// summed in 64-bit integer arithmetic before dividing, so values beyond 2^53
// do not lose precision to intermediate float64 rounding.
// It returns math.NaN() if the collection is empty.
//
// AverageBy is a generic method: the numeric type TNumber is inferred from the
// selector function.
func (q Query[T]) AverageBy[TNumber Number](selector func(T) TNumber) float64 {
	return averageBy(q, selector)
}

// Max returns the maximum value in a collection of ordered values and a
// boolean reporting whether the collection was non-empty.
func Max[T cmp.Ordered](q Query[T]) (T, bool) {
	var r T
	found := false
	q.Iterate(func(item T) bool {
		if !found || cmp.Compare(item, r) > 0 {
			r = item
			found = true
		}
		return true
	})
	return r, found
}

// Min returns the minimum value in a collection of ordered values and a
// boolean reporting whether the collection was non-empty.
func Min[T cmp.Ordered](q Query[T]) (T, bool) {
	var r T
	found := false
	q.Iterate(func(item T) bool {
		if !found || cmp.Compare(item, r) < 0 {
			r = item
			found = true
		}
		return true
	})
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
	q.Iterate(func(item T) bool {
		key := selector(item)
		if !found || cmp.Compare(key, rKey) > 0 {
			r, rKey = item, key
			found = true
		}
		return true
	})
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
	q.Iterate(func(item T) bool {
		key := selector(item)
		if !found || cmp.Compare(key, rKey) < 0 {
			r, rKey = item, key
			found = true
		}
		return true
	})
	return r, found
}
