package linq

import "reflect"

// Iterator is an alias for function to iterate over data.
type Iterator func() (item interface{}, ok bool)

// Query is the type returned from query functions. It can be iterated manually
// as shown in the example.
type Query struct {
	Iterate func() Iterator
}

// KeyValue is a type that is used to iterate over a map (if query is created
// from a map). This type is also used by ToMap() method to output result of a
// query into a map.
type KeyValue struct {
	Key   interface{}
	Value interface{}
}

// Iterable is an interface that has to be implemented by a custom collection in
// order to work with linq.
type Iterable interface {
	Iterate() Iterator
}

// From initializes a linq query with passed slice, array or map as the source.
// String, channel or struct implementing Iterable interface can be used as an
// input. In this case From delegates it to FromString, FromChannel and
// FromIterable internally.
func From(source interface{}) Query {
	src := reflect.ValueOf(source)

	switch src.Kind() {
	case reflect.Slice, reflect.Array:
		len := src.Len()

		return Query{
			Iterate: func() Iterator {
				index := 0

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						item = src.Index(index).Interface()
						index++
					}

					return
				}
			},
		}
	case reflect.Map:
		len := src.Len()

		return Query{
			Iterate: func() Iterator {
				index := 0
				keys := src.MapKeys()

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						key := keys[index]
						item = KeyValue{
							Key:   key.Interface(),
							Value: src.MapIndex(key).Interface(),
						}

						index++
					}

					return
				}
			},
		}
	case reflect.String:
		return FromString(source.(string))
	case reflect.Chan:
		if _, ok := source.(chan interface{}); ok {
			return FromChannel(source.(chan interface{}))
		} else {
			return FromChannelT(source)
		}
	default:
		return FromIterable(source.(Iterable))
	}
}

// FromChannel initializes a linq query with passed channel, linq iterates over
// channel until it is closed.
func FromChannel(source <-chan interface{}) Query {
	return Query{
		Iterate: func() Iterator {
			return func() (item interface{}, ok bool) {
				item, ok = <-source
				return
			}
		},
	}
}

// FromChannelT is the typed version of FromChannel.
//
//   - source is of type "chan TSource"
//
// NOTE: FromChannel has better performance than FromChannelT.
func FromChannelT(source interface{}) Query {
	src := reflect.ValueOf(source)
	return Query{
		Iterate: func() Iterator {
			return func() (interface{}, bool) {
				value, ok := src.Recv()
				return value.Interface(), ok
			}
		},
	}
}

// FromString initializes a linq query with passed string, linq iterates over
// runes of string.
func FromString(source string) Query {
	runes := []rune(source)
	len := len(runes)

	return Query{
		Iterate: func() Iterator {
			index := 0

			return func() (item interface{}, ok bool) {
				ok = index < len
				if ok {
					item = runes[index]
					index++
				}

				return
			}
		},
	}
}

// FromIterable initializes a linq query with custom collection passed. This
// collection has to implement Iterable interface, linq iterates over items,
// that has to implement Comparable interface or be basic types.
func FromIterable(source Iterable) Query {
	return Query{
		Iterate: source.Iterate,
	}
}

// Range generates a sequence of integral numbers within a specified range.
func Range(start, count int) Query {
	return Query{
		Iterate: func() Iterator {
			index := 0
			current := start

			return func() (item interface{}, ok bool) {
				if index >= count {
					return nil, false
				}

				item, ok = current, true

				index++
				current++
				return
			}
		},
	}
}

// Repeat generates a sequence that contains one repeated value.
func Repeat(value interface{}, count int) Query {
	return Query{
		Iterate: func() Iterator {
			index := 0

			return func() (item interface{}, ok bool) {
				if index >= count {
					return nil, false
				}

				item, ok = value, true

				index++
				return
			}
		},
	}
}
