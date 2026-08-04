package linq

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func (q Query[T]) SumBy[N number](selector func(T) N) N {
	var result N
	q.iterate(func(value T) bool {
		result += selector(value)
		return true
	})
	return result
}
