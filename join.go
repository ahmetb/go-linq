package linq

import "reflect"

// Join correlates the elements of two collections based on matching keys.
//
// A join refers to the operation of correlating the elements of two sources of
// information based on a common key. Join brings the two information sources
// and the keys by which they are matched together in one method call. This
// differs from the use of SelectMany, which requires more than one method call
// to perform the same operation.
//
// Join is a generic method: the inner element type TInner, the key type TKey, and the
// result type TResult are all inferred from the supplied functions. The key type TKey
// must be comparable.
//
// Join preserves the order of the elements of outer collection, and for each of
// these elements, the order of the matching elements of inner.
//
// Elements whose key is nil take no part in the join: like .NET LINQ, Join
// never matches a nil key, not even against another nil key.
func (q Query[T]) Join[TInner any, TKey comparable, TResult any](inner Query[TInner],
	outerKeySelector func(T) TKey,
	innerKeySelector func(TInner) TKey,
	resultSelector func(outer T, inner TInner) TResult) Query[TResult] {

	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			innerLookup := buildJoinLookup(inner, innerKeySelector)
			interfaceKey := reflect.TypeFor[TKey]().Kind() == reflect.Interface

			q.Iterate(func(outerItem T) bool {
				for _, innerItem := range joinGroupFor(innerLookup,
					outerKeySelector(outerItem), interfaceKey) {
					result := resultSelector(outerItem, innerItem)
					if !yield(result) {
						return false
					}
				}
				return true
			})
		},
	}
}

// LeftJoin correlates two collections by key while retaining every element of
// the outer collection. When an outer element has no match — a nil key never
// matches, not even another nil key — resultSelector is called with the zero
// value of TInner, which it cannot tell apart from a matched element that is
// itself the zero value.
//
// Inner is indexed from the first outer element, so an empty outer never reads
// inner. An empty inner has no such shortcut: every outer element still owes a
// result.
func (q Query[T]) LeftJoin[TInner any, TKey comparable, TResult any](inner Query[TInner],
	outerKeySelector func(T) TKey,
	innerKeySelector func(TInner) TKey,
	resultSelector func(outer T, inner TInner) TResult) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			var innerLookup map[TKey][]TInner
			interfaceKey := reflect.TypeFor[TKey]().Kind() == reflect.Interface

			q.Iterate(func(outerItem T) bool {
				if innerLookup == nil {
					innerLookup = buildJoinLookup(inner, innerKeySelector)
				}

				innerGroup := joinGroupFor(innerLookup, outerKeySelector(outerItem), interfaceKey)
				if len(innerGroup) == 0 {
					var zero TInner
					return yield(resultSelector(outerItem, zero))
				}

				for _, innerItem := range innerGroup {
					if !yield(resultSelector(outerItem, innerItem)) {
						return false
					}
				}
				return true
			})
		},
	}
}

// RightJoin correlates two collections by key while retaining every element
// of the inner collection. When an inner element has no match — a nil key never
// matches, not even another nil key — resultSelector is called with the zero
// value of T, which it cannot tell apart from a matched element that is itself
// the zero value.
//
// Outer is indexed from the first inner element, so an empty inner never reads
// outer. An empty outer has no such shortcut: every inner element still owes a
// result.
func (q Query[T]) RightJoin[TInner any, TKey comparable, TResult any](inner Query[TInner],
	outerKeySelector func(T) TKey,
	innerKeySelector func(TInner) TKey,
	resultSelector func(outer T, inner TInner) TResult) Query[TResult] {
	// Keep the mirrored implementation to avoid an argument-swapping call per result.
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			var outerLookup map[TKey][]T
			interfaceKey := reflect.TypeFor[TKey]().Kind() == reflect.Interface

			inner.Iterate(func(innerItem TInner) bool {
				if outerLookup == nil {
					outerLookup = buildJoinLookup(q, outerKeySelector)
				}

				outerGroup := joinGroupFor(outerLookup, innerKeySelector(innerItem), interfaceKey)
				if len(outerGroup) == 0 {
					var zero T
					return yield(resultSelector(zero, innerItem))
				}

				for _, outerItem := range outerGroup {
					if !yield(resultSelector(outerItem, innerItem)) {
						return false
					}
				}
				return true
			})
		},
	}
}

// FullJoin correlates two collections by key while retaining every element of
// both collections. When an element has no match — a nil key never matches,
// not even another nil key — resultSelector is called with the zero value for
// the missing side, which it cannot tell apart from a matched element that is
// itself the zero value.
//
// Inner is indexed before outer is read. Results for outer elements come
// first, followed by unmatched inner groups in first-seen key order.
func (q Query[T]) FullJoin[TInner any, TKey comparable, TResult any](inner Query[TInner],
	outerKeySelector func(T) TKey,
	innerKeySelector func(TInner) TKey,
	resultSelector func(outer T, inner TInner) TResult) Query[TResult] {
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			innerGroups, innerLookup := buildFullJoinGroups(inner, innerKeySelector)
			matched := make([]bool, len(innerGroups))
			interfaceKey := reflect.TypeFor[TKey]().Kind() == reflect.Interface

			noMatch := make([]TInner, 1)

			for outerItem := range q.Iterate {
				innerItems := noMatch
				outerKey := outerKeySelector(outerItem)
				if !interfaceKey || !isNilInterfaceKey(outerKey) {
					if index, ok := innerLookup[outerKey]; ok {
						matched[index] = true
						innerItems = innerGroups[index]
					}
				}

				for _, innerItem := range innerItems {
					if !yield(resultSelector(outerItem, innerItem)) {
						return
					}
				}
			}

			var zero T
			for i, innerGroup := range innerGroups {
				if matched[i] {
					continue
				}
				for _, innerItem := range innerGroup {
					if !yield(resultSelector(zero, innerItem)) {
						return
					}
				}
			}
		},
	}
}

// isNilInterfaceKey reports whether an interface key is nil or holds a typed nil.
func isNilInterfaceKey(key any) bool {
	value := reflect.ValueOf(key)
	if !value.IsValid() {
		return true
	}
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	case reflect.UnsafePointer:
		return value.IsZero()
	}
	return false
}

// joinGroupFor returns nil for interface keys that hold a typed nil, which
// may be unhashable and must not reach a map lookup.
func joinGroupFor[T any, TKey comparable](lookup map[TKey][]T, key TKey, interfaceKey bool) []T {
	if interfaceKey && isNilInterfaceKey(key) {
		return nil
	}
	return lookup[key]
}

// buildJoinLookup indexes non-nil keys. .NET LINQ joins do not match null keys.
func buildJoinLookup[T any, TKey comparable](source Query[T], keySelector func(T) TKey) map[TKey][]T {
	lookup := make(map[TKey][]T)
	switch reflect.TypeFor[TKey]().Kind() {
	case reflect.Chan, reflect.Pointer, reflect.UnsafePointer:
		// Nilable statically, so the zero value is the only nil there is.
		var zero TKey
		source.Iterate(func(item T) bool {
			key := keySelector(item)
			if key != zero {
				lookup[key] = append(lookup[key], item)
			}
			return true
		})

	case reflect.Interface:
		// A typed nil pointer in an interface is not the nil interface, so
		// nothing short of the dynamic value can tell them apart. This is the
		// one path that pays for reflection per element.
		source.Iterate(func(item T) bool {
			key := keySelector(item)
			if isNilInterfaceKey(key) {
				return true
			}
			lookup[key] = append(lookup[key], item)
			return true
		})

	default:
		// No value of TKey can be nil.
		source.Iterate(func(item T) bool {
			key := keySelector(item)
			lookup[key] = append(lookup[key], item)
			return true
		})
	}
	return lookup
}

// buildFullJoinGroups groups elements in first-seen key order. Nil keys share
// one unmatched group and stay out of the index.
func buildFullJoinGroups[T any, TKey comparable](source Query[T],
	keySelector func(T) TKey) ([][]T, map[TKey]int) {
	keyKind := reflect.TypeFor[TKey]().Kind()
	zeroIsNil := keyKind == reflect.Chan || keyKind == reflect.Pointer || keyKind == reflect.UnsafePointer
	interfaceKey := keyKind == reflect.Interface
	var zeroKey TKey

	var groups [][]T
	index := make(map[TKey]int)
	nilGroupIndex := -1
	source.Iterate(func(item T) bool {
		key := keySelector(item)
		if (zeroIsNil && key == zeroKey) || (interfaceKey && isNilInterfaceKey(key)) {
			if nilGroupIndex < 0 {
				nilGroupIndex = len(groups)
				groups = append(groups, nil)
			}
			groups[nilGroupIndex] = append(groups[nilGroupIndex], item)
			return true
		}

		i, ok := index[key]
		if !ok {
			i = len(groups)
			index[key] = i
			groups = append(groups, nil)
		}
		groups[i] = append(groups[i], item)
		return true
	})
	return groups, index
}
