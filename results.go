package linq

func (q Query[T]) Results() []T {
	return q.toSlice()
}
