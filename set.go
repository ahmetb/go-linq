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

// setOps bundles the element-type-dependent operations of the set operators:
// constructing a seenSet and comparing two elements for equality.
type setOps[T any] struct {
	newSet func() seenSet[T]
	equal  func(T, T) bool
}

func typedOps[T comparable]() any {
	return setOps[T]{
		newSet: typedSet[T],
		equal:  func(a, b T) bool { return a == b },
	}
}

// opsFor returns operations backed by strongly-typed maps and direct ==
// comparison when T is a basic comparable kind, so that no per-element boxing
// occurs. For all other element types it falls back to operating on boxed
// (interface) values, which panics at runtime if T is not comparable
// (matching v4 behavior).
//
// The type-switch cases are exact: a named type whose underlying type is a
// basic kind (e.g. "type ID int") takes the fallback path, preserving its own
// equality semantics through boxing.
func opsFor[T any]() setOps[T] {
	var o any
	switch any(*new(T)).(type) {
	case int:
		o = typedOps[int]()
	case int8:
		o = typedOps[int8]()
	case int16:
		o = typedOps[int16]()
	case int32:
		o = typedOps[int32]()
	case int64:
		o = typedOps[int64]()
	case uint:
		o = typedOps[uint]()
	case uint8:
		o = typedOps[uint8]()
	case uint16:
		o = typedOps[uint16]()
	case uint32:
		o = typedOps[uint32]()
	case uint64:
		o = typedOps[uint64]()
	case uintptr:
		o = typedOps[uintptr]()
	case float32:
		o = typedOps[float32]()
	case float64:
		o = typedOps[float64]()
	case complex64:
		o = typedOps[complex64]()
	case complex128:
		o = typedOps[complex128]()
	case string:
		o = typedOps[string]()
	case bool:
		o = typedOps[bool]()
	default:
		return setOps[T]{
			newSet: boxedSet[T],
			equal:  func(a, b T) bool { return any(a) == any(b) },
		}
	}
	// Never fails: in each case above T is that exact concrete type, so
	// setOps[T] is the same instantiated type as the stored value.
	return o.(setOps[T])
}

func newSeenSet[T any]() seenSet[T] { return opsFor[T]().newSet() }

func equalFor[T any]() func(T, T) bool { return opsFor[T]().equal }
