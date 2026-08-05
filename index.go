package linq

// IndexOf searches for an element that matches the conditions defined by a specified predicate
// and returns the zero-based index of the first occurrence within the collection. This method
// returns -1 if an item that matches the conditions is not found.
func (q Query[T]) IndexOf(predicate func(T) bool) int {
	result := -1
	index := 0
	q.Iterate(func(item T) bool {
		if predicate(item) {
			result = index
			return false
		}
		index++
		return true
	})
	return result
}
