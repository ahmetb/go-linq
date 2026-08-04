package linq

// ContainsBy reports whether any value projects to key.
func (q Query[T]) ContainsBy[K comparable](key K, selector func(T) K) bool {
	found := false
	q.iterate(func(value T) bool {
		if selector(value) == key {
			found = true
			return false
		}
		return true
	})
	return found
}
