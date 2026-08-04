package linq

import "math"

// AverageBy returns the arithmetic mean of selected numeric values, or NaN when empty.
func (q Query[T]) AverageBy[N Number](selector func(T) N) float64 {
	sum := 0.0
	count := 0
	q.iterate(func(value T) bool {
		sum += float64(selector(value))
		count++
		return true
	})
	if count == 0 {
		return math.NaN()
	}
	return sum / float64(count)
}
