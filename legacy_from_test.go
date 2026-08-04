package linq

import (
	"context"
	"fmt"
	"iter"
	"reflect"
)

// Query is the type returned from query functions. It can be iterated manually
// as shown in the example.
type legacyQuery struct {
	Iterate iter.Seq[any]
}

// KeyValue is a type used to iterate over a map. This type is also used by ToMap()
// method to output the result of a query into a map.
type legacyKeyValue struct {
	Key   any
	Value any
}

// Iterable is an interface that has to be implemented by a custom collection
// to work with linq.
type legacyIterable interface {
	Iterate() iter.Seq[any]
}

// FromSlice initializes a linq query with a passed slice.
func legacyFromSlice[S ~[]T, T any](source S) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// FromMap initializes a linq query with a passed map.
func legacyFromMap[M ~map[K]V, K comparable, V any](source M) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			for k, v := range source {
				if !yield(legacyKeyValue{
					Key:   k,
					Value: v,
				}) {
					return
				}
			}
		},
	}
}

// FromChannel initializes a linq query with a passed channel, linq iterates over
// the channel until it is closed.
func legacyFromChannel[T any](source <-chan T) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			for item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// FromChannelWithContext initializes a linq query with a passed channel
// and stops iterating either when the channel is closed or when the context is canceled.
func legacyFromChannelWithContext[T any](ctx context.Context, source <-chan T) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			for {
				select {
				case <-ctx.Done():
					// Context canceled or deadline exceeded
					return
				case item, ok := <-source:
					if !ok || !yield(item) {
						// Channel closed or Consumer stopped early
						return
					}
				}
			}
		},
	}
}

// FromString initializes a query from a string, iterating over its runes.
func legacyFromString[S ~string](source S) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			for _, ch := range string(source) {
				if !yield(ch) {
					return
				}
			}
		},
	}
}

// FromIterable initializes a linq query with a custom collection passed. This
// collection has to implement Iterable.
func legacyFromIterable(source legacyIterable) legacyQuery {
	return legacyQuery{
		Iterate: source.Iterate(),
	}
}

// From initializes a Query from a supported data source by inspecting its
// type at runtime. It panics if the source type is not supported.
//
// NOTE: It is recommended to call the specific From* function directly
// (e.g., FromSlice, FromMap, etc.). This unified function is less efficient
// because it relies on runtime reflection.
func legacyFrom(source any) legacyQuery {
	if source == nil {
		return legacyQuery{
			Iterate: func(yield func(any) bool) {},
		}
	}

	switch s := source.(type) {
	case string:
		return legacyFromString(s)
	case legacyIterable:
		return legacyFromIterable(s)
	}

	sourceValue := reflect.ValueOf(source)
	switch sourceValue.Kind() {
	case reflect.Slice, reflect.Array:
		return legacyQuery{
			Iterate: func(yield func(any) bool) {
				length := sourceValue.Len()
				for i := 0; i < length; i++ {
					if !yield(sourceValue.Index(i).Interface()) {
						return
					}
				}
			},
		}

	case reflect.Map:
		return legacyQuery{
			Iterate: func(yield func(any) bool) {
				for _, key := range sourceValue.MapKeys() {
					value := sourceValue.MapIndex(key)
					if !yield(legacyKeyValue{Key: key.Interface(), Value: value.Interface()}) {
						return
					}
				}
			},
		}

	case reflect.Chan:
		return legacyQuery{
			Iterate: func(yield func(any) bool) {
				for {
					value, ok := sourceValue.Recv()
					if !ok || !yield(value.Interface()) {
						return
					}
				}
			},
		}

	default:
		panic(fmt.Sprintf("unsupported type for From: %T", source))
	}
}

// Range generates a sequence of integral numbers within a specified range.
func legacyRange(start, count int) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			end := start + count
			for i := start; i < end; i++ {
				if !yield(i) {
					return
				}
			}
		},
	}
}

// Repeat generates a sequence that contains one repeated value.
func legacyRepeat[T any](value T, count int) legacyQuery {
	return legacyQuery{
		Iterate: func(yield func(any) bool) {
			for i := 0; i < count; i++ {
				if !yield(value) {
					return
				}
			}
		},
	}
}
