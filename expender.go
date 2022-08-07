package linq

var _ Expender[int] = &expender[int, int]{}
var _ Expended[int, int] = &expender[int, int]{}
var _ Expended3[int, int, int] = &expender3[int, int, int]{}

func (q QueryG[T]) Expend(e Expender[T]) Expender[T] {
	e.Expend(q)
	return e
}

type Expender[T any] interface {
	Expend(q QueryG[T]) any
}

type Expended[T1, T2 any] interface {
	Select(selector func(T1) T2) QueryG[T2]
	SelectIndexed(selector func(int, T1) T2) QueryG[T2]
	SelectMany(selector func(T1) QueryG[T2]) QueryG[T2]
	SelectManyIndexed(selector func(int, T1) QueryG[T2]) QueryG[T2]
}

type Expended3[T1, T2, T3 any] interface {
	Zip(q2 QueryG[T2],
		resultSelector func(T1, T2) T3) QueryG[T3]
	SelectManyBy(selector func(T1) QueryG[T2],
		resultSelector func(T2, T1) T3) QueryG[T3]
	SelectManyByIndexed(selector func(int, T1) QueryG[T2],
		resultSelector func(T2, T1) T3) QueryG[T3]
}

type expender[T1, T2 any] struct {
	q QueryG[T1]
}

func (e *expender[T1, T2]) Expend(q QueryG[T1]) any {
	e.q = q
	return e
}

func To2[T1, T2 any]() Expender[T1] {
	return &expender[T1, T2]{}
}

type expender3[T1, T2, T3 any] struct {
	q QueryG[T1]
}

func (e *expender3[T1, T2, T3]) Expend(q QueryG[T1]) any {
	e.q = q
	return e
}

func To3[T1, T2, T3 any]() Expender[T1] {
	return &expender3[T1, T2, T3]{}
}
