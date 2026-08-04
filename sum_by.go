package linq

// Number contains the built-in integer and floating-point types and their named forms.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// SumBy returns the sum of selected numeric values, or zero when empty.
func (q Query[T]) SumBy[N Number](selector func(T) N) N {
	var result N
	q.iterate(func(value T) bool {
		result += selector(value)
		return true
	})
	return result
}
