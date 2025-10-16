package linq

import (
	"context"
	"fmt"
	"iter"
	"reflect"
)

// Query is the type returned from query functions. It can be iterated manually
// as shown in the example.
type Query struct {
	Iterate iter.Seq[any]
}

// KeyValue is a type used to iterate over a map. This type is also used by ToMap()
// method to output the result of a query into a map.
type KeyValue struct {
	Key   any
	Value any
}

// Iterable is an interface that has to be implemented by a custom collection
// to work with linq.
type Iterable interface {
	Iterate() iter.Seq[any]
}

// FromSlice initializes a linq query with a passed slice.
func FromSlice[S ~[]T, T any](source S) Query {
	return Query{
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
func FromMap[M ~map[K]V, K comparable, V any](source M) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			for k, v := range source {
				if !yield(KeyValue{
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
func FromChannel[T any](source <-chan T) Query {
	return Query{
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
func FromChannelWithContext[T any](ctx context.Context, source <-chan T) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			for {
				select {
				case <-ctx.Done():
					// Context canceled or deadline exceeded
					return
				case item, ok := <-source:
					if !ok {
						// Channel closed
						return
					}
					if !yield(item) {
						// Consumer stopped early
						return
					}
				}
			}
		},
	}
}

// FromString initializes a query from a string, iterating over its runes.
func FromString[S ~string](source S) Query {
	return Query{
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
func FromIterable(source Iterable) Query {
	return Query{
		Iterate: source.Iterate(),
	}
}

// From initializes a Query from a supported data source by inspecting its
// type at runtime. It panics if the source type is not supported.
//
// NOTE: It is recommended to call the specific From* function directly
// (e.g., FromSlice, FromMap, etc.). This unified function is less efficient
// because it relies on runtime reflection.
func From(source any) Query {
	if source == nil {
		return Query{
			Iterate: func(yield func(any) bool) {},
		}
	}

	switch s := source.(type) {
	case string:
		return FromString(s)
	case Iterable:
		return FromIterable(s)
	}

	sourceValue := reflect.ValueOf(source)
	switch sourceValue.Kind() {
	case reflect.Slice, reflect.Array:
		return Query{
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
		return Query{
			Iterate: func(yield func(any) bool) {
				for _, key := range sourceValue.MapKeys() {
					value := sourceValue.MapIndex(key)
					if !yield(KeyValue{Key: key.Interface(), Value: value.Interface()}) {
						return
					}
				}
			},
		}

	case reflect.Chan:
		return Query{
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
func Range(start, count int) Query {
	return Query{
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
func Repeat[T any](value T, count int) Query {
	return Query{
		Iterate: func(yield func(any) bool) {
			for i := 0; i < count; i++ {
				if !yield(value) {
					return
				}
			}
		},
	}
}
