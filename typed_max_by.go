package linq

import "cmp"

func (q query[T]) MaxBy[K cmp.Ordered](selector func(T) K) T {
	var result T
	var maximum K
	found := false
	q.iterate(func(value T) bool {
		key := selector(value)
		if !found || cmp.Compare(key, maximum) > 0 {
			result = value
			maximum = key
			found = true
		}
		return true
	})
	return result
}
