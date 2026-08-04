package linq

import "cmp"

func (q Query[T]) MinBy[K cmp.Ordered](selector func(T) K) T {
	var result T
	var minimum K
	found := false
	q.iterate(func(value T) bool {
		key := selector(value)
		if !found || cmp.Compare(key, minimum) < 0 {
			result = value
			minimum = key
			found = true
		}
		return true
	})
	return result
}
