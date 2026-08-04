package linq

type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

func fromMap[M ~map[K]V, K comparable, V any](source M) Query[KeyValue[K, V]] {
	return Query[KeyValue[K, V]]{iterate: func(yield func(KeyValue[K, V]) bool) {
		for key, value := range source {
			if !yield(KeyValue[K, V]{Key: key, Value: value}) {
				return
			}
		}
	}}
}
