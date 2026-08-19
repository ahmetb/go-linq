package linq

// Position represents a position counted from the start or end of a sequence.
// Its zero value is PositionFromStart(0).
// PositionFromEnd(0) points just past the sequence and does not identify an element.
type Position struct {
	value   int
	fromEnd bool
}

// NewPosition returns a Position counted from the end of a sequence when fromEnd is
// true, or from the start otherwise. It panics if value is negative.
func NewPosition(value int, fromEnd bool) Position {
	if value < 0 {
		panic("linq: position must be non-negative")
	}
	return Position{value: value, fromEnd: fromEnd}
}

// PositionFromStart returns a Position counted from the start of a sequence. It
// panics if value is negative.
func PositionFromStart(value int) Position {
	return NewPosition(value, false)
}

// PositionFromEnd returns a Position counted from the end of a sequence.
// PositionFromEnd(1) identifies the last element; PositionFromEnd(0) identifies the
// position immediately after it. It panics if value is negative.
func PositionFromEnd(value int) Position {
	return NewPosition(value, true)
}

func (p Position) offset(length int) int {
	if p.fromEnd {
		return max(length-p.value, 0)
	}
	return min(p.value, length)
}

// ElementAt returns the element at position and a boolean reporting whether that
// position exists.
func (q Query[T]) ElementAt(position Position) (T, bool) {
	value := position.value
	if position.fromEnd {
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
	current := 0
	q.Iterate(func(item T) bool {
		if current == value {
			result = item
			found = true
			return false
		}
		current++
		return true
	})
	return result, found
}

// Index pairs each item with its zero-based index as a KeyValue.
func Index[T any](q Query[T]) Query[KeyValue[int, T]] {
	return q.SelectIndexed(func(index int, item T) KeyValue[int, T] {
		return KeyValue[int, T]{Key: index, Value: item}
	})
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
