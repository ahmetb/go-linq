package linq

func (q query[T]) Results() []T {
	return q.toSlice()
}
