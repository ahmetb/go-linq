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

			q.Iterate(func(outerItem T) bool {
				outerKey := outerKeySelector(outerItem)

				if innerGroup, ok := innerLookup[outerKey]; ok {
					for _, innerItem := range innerGroup {
						result := resultSelector(outerItem, innerItem)
						if !yield(result) {
							return false
						}
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
// value of TInner.
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

			q.Iterate(func(outerItem T) bool {
				if innerLookup == nil {
					innerLookup = buildJoinLookup(inner, innerKeySelector)
				}

				innerGroup, ok := innerLookup[outerKeySelector(outerItem)]
				if !ok {
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
// value of T.
//
// Outer is indexed from the first inner element, so an empty inner never reads
// outer. An empty outer has no such shortcut: every inner element still owes a
// result.
func (q Query[T]) RightJoin[TInner any, TKey comparable, TResult any](inner Query[TInner],
	outerKeySelector func(T) TKey,
	innerKeySelector func(TInner) TKey,
	resultSelector func(outer T, inner TInner) TResult) Query[TResult] {
	// ponytail: Delegating to LeftJoin is shorter but adds an argument-swapping
	// call per result (~10% in paired benchmarks); collapse if that call becomes free.
	return Query[TResult]{
		Iterate: func(yield func(TResult) bool) {
			var outerLookup map[TKey][]T

			inner.Iterate(func(innerItem TInner) bool {
				if outerLookup == nil {
					outerLookup = buildJoinLookup(q, outerKeySelector)
				}

				outerGroup, ok := outerLookup[innerKeySelector(innerItem)]
				if !ok {
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
			value := reflect.ValueOf(key)
			if !value.IsValid() {
				return true
			}
			switch value.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice:
				if value.IsNil() {
					return true
				}
			case reflect.UnsafePointer:
				// Not IsNil: reflect panics on it for this kind.
				if value.IsZero() {
					return true
				}
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
