package linq

// seenSet tracks values of type T for the set operators that have no key
// selector (Distinct, Union, Except, Intersect). add inserts a value and
// reports whether it was absent; has reports membership without inserting;
// del removes a value.
type seenSet[T any] struct {
	add func(T) bool
	has func(T) bool
	del func(T)
}

func typedSet[T comparable]() seenSet[T] {
	m := make(map[T]struct{})
	return seenSet[T]{
		add: func(v T) bool {
			if _, ok := m[v]; ok {
				return false
			}
			m[v] = struct{}{}
			return true
		},
		has: func(v T) bool {
			_, ok := m[v]
			return ok
		},
		del: func(v T) {
			delete(m, v)
		},
	}
}

func boxedSet[T any]() seenSet[T] {
	m := make(map[any]struct{})
	return seenSet[T]{
		add: func(v T) bool {
			if _, ok := m[v]; ok {
				return false
			}
			m[v] = struct{}{}
			return true
		},
		has: func(v T) bool {
			_, ok := m[v]
			return ok
		},
		del: func(v T) {
			delete(m, v)
		},
	}
}

// newSeenSet returns a set of T backed by a strongly-typed map when T is a
// basic comparable kind, so that no per-element boxing occurs. For all other
// element types it falls back to a map keyed by boxed values, which panics at
// runtime if T is not comparable (matching v4 behavior).
//
// The type assertions below are exact: a named type whose underlying type is
// a basic kind (e.g. "type ID int") takes the fallback path, preserving its
// own equality semantics through boxing.
func newSeenSet[T any]() seenSet[T] {
	var s any
	switch any(*new(T)).(type) {
	case int:
		s = typedSet[int]()
	case int8:
		s = typedSet[int8]()
	case int16:
		s = typedSet[int16]()
	case int32:
		s = typedSet[int32]()
	case int64:
		s = typedSet[int64]()
	case uint:
		s = typedSet[uint]()
	case uint8:
		s = typedSet[uint8]()
	case uint16:
		s = typedSet[uint16]()
	case uint32:
		s = typedSet[uint32]()
	case uint64:
		s = typedSet[uint64]()
	case uintptr:
		s = typedSet[uintptr]()
	case float32:
		s = typedSet[float32]()
	case float64:
		s = typedSet[float64]()
	case complex64:
		s = typedSet[complex64]()
	case complex128:
		s = typedSet[complex128]()
	case string:
		s = typedSet[string]()
	case bool:
		s = typedSet[bool]()
	default:
		return boxedSet[T]()
	}
	return s.(seenSet[T])
}

func typedEqual[T comparable]() any {
	return func(a, b T) bool { return a == b }
}

// equalFor returns an equality function for T: a direct == comparison when T
// is a basic comparable kind, and comparison of boxed values otherwise. The
// boxed comparison panics at runtime if T is not comparable (matching v4
// behavior).
func equalFor[T any]() func(T, T) bool {
	var f any
	switch any(*new(T)).(type) {
	case int:
		f = typedEqual[int]()
	case int8:
		f = typedEqual[int8]()
	case int16:
		f = typedEqual[int16]()
	case int32:
		f = typedEqual[int32]()
	case int64:
		f = typedEqual[int64]()
	case uint:
		f = typedEqual[uint]()
	case uint8:
		f = typedEqual[uint8]()
	case uint16:
		f = typedEqual[uint16]()
	case uint32:
		f = typedEqual[uint32]()
	case uint64:
		f = typedEqual[uint64]()
	case uintptr:
		f = typedEqual[uintptr]()
	case float32:
		f = typedEqual[float32]()
	case float64:
		f = typedEqual[float64]()
	case complex64:
		f = typedEqual[complex64]()
	case complex128:
		f = typedEqual[complex128]()
	case string:
		f = typedEqual[string]()
	case bool:
		f = typedEqual[bool]()
	default:
		return func(a, b T) bool { return any(a) == any(b) }
	}
	return f.(func(T, T) bool)
}
