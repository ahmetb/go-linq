package linq

type keyValue[K comparable, V any] struct {
	Key   K
	Value V
}

func fromMap[M ~map[K]V, K comparable, V any](source M) Query[keyValue[K, V]] {
	return Query[keyValue[K, V]]{iterate: func(yield func(keyValue[K, V]) bool) {
		for key, value := range source {
			if !yield(keyValue[K, V]{Key: key, Value: value}) {
				return
			}
		}
	}}
}
