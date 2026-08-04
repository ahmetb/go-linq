package linq

import "context"

func fromChannelWithContext[T any](ctx context.Context, source <-chan T) query[T] {
	return query[T]{iterate: func(yield func(T) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-source:
				if !ok || !yield(value) {
					return
				}
			}
		}
	}}
}
