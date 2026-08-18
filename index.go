package linq

// Index represents a position counted from the start or end of a sequence.
// Its zero value is IndexFromStart(0).
type Index struct {
	value   int
	fromEnd bool
}

// NewIndex returns an Index counted from the end of a sequence when fromEnd is
// true, or from the start otherwise. It panics if value is negative.
func NewIndex(value int, fromEnd bool) Index {
	if value < 0 {
		panic("linq: index must be non-negative")
	}
	return Index{value: value, fromEnd: fromEnd}
}

// IndexFromStart returns an Index counted from the start of a sequence. It
// panics if value is negative.
func IndexFromStart(value int) Index {
	return NewIndex(value, false)
}

// IndexFromEnd returns an Index counted from the end of a sequence.
// IndexFromEnd(1) identifies the last element; IndexFromEnd(0) identifies the
// position immediately after it. It panics if value is negative.
func IndexFromEnd(value int) Index {
	return NewIndex(value, true)
}

func (i Index) offset(length int) int {
	if i.fromEnd {
		return max(length-i.value, 0)
	}
	return min(i.value, length)
}

// ElementAt returns the element at index and a boolean reporting whether that
// position exists.
func (q Query[T]) ElementAt(index Index) (T, bool) {
	value := index.value
	if index.fromEnd {
		var zero T
		if value == 0 {
			return zero, false
		}

		items, _ := lastItems(q, value)
		if len(items) < value {
			return zero, false
		}
		return items[0], true
	}

	var result T
	found := false
	position := 0
	q.Iterate(func(item T) bool {
		if position == value {
			result = item
			found = true
			return false
		}
		position++
		return true
	})
	return result, found
}

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
