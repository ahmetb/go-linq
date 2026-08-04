package linq

import "iter"

type typedIterable[T any] interface {
	Iterate() iter.Seq[T]
}

func fromIterable[T any](source typedIterable[T]) query[T] {
	return query[T]{iterate: source.Iterate()}
}
