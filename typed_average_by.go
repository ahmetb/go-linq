package linq

import "math"

func (q query[T]) AverageBy[N number](selector func(T) N) float64 {
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
