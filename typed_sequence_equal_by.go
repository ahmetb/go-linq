package linq

import "iter"

func (q query[T]) SequenceEqualBy[U any, K comparable](
	other query[U],
	leftSelector func(T) K,
	rightSelector func(U) K,
) bool {
	nextRight, stopRight := iter.Pull(other.iterate)
	defer stopRight()

	equal := true
	q.iterate(func(left T) bool {
		right, ok := nextRight()
		if !ok || leftSelector(left) != rightSelector(right) {
			equal = false
			return false
		}
		return true
	})
	if !equal {
		return false
	}

	_, rightOK := nextRight()
	return !rightOK
}
