package linq

import "iter"

type typedIterable[T any] interface {
	Iterate() iter.Seq[T]
}

func fromIterable[T any](source typedIterable[T]) Query[T] {
	return Query[T]{iterate: source.Iterate()}
}
