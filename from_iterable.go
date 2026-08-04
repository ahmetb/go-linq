package linq

import "iter"

// Iterable is implemented by custom collections that expose a typed sequence.
type Iterable[T any] interface {
	Iterate() iter.Seq[T]
}

func fromIterable[T any](source Iterable[T]) Query[T] {
	return Query[T]{iterate: source.Iterate()}
}
