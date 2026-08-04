package linq

// Group is a type used to store the result of the GroupBy method.
type Group[TKey comparable, TElement any] struct {
	Key   TKey
	Group []TElement
}

// GroupBy method groups the elements of a collection according to a specified
// key selector function and projects the elements for each group by using a
// specified function.
//
// GroupBy is a generic method: the key type TKey and the element type TElement are
// inferred from the supplied functions. The key type TKey must be comparable.
//
// Groups are yielded in the order of the first appearance of their key in the
// source collection, and elements within each group preserve the order they
// appear in the source.
func (q Query[T]) GroupBy[TKey comparable, TElement any](keySelector func(T) TKey,
	elementSelector func(T) TElement) Query[Group[TKey, TElement]] {
	return Query[Group[TKey, TElement]]{
		Iterate: func(yield func(Group[TKey, TElement]) bool) {
			groups := make(map[TKey][]TElement)
			var keys []TKey

			for item := range q.Iterate {
				key := keySelector(item)
				if _, ok := groups[key]; !ok {
					keys = append(keys, key)
				}
				groups[key] = append(groups[key], elementSelector(item))
			}

			for _, key := range keys {
				group := Group[TKey, TElement]{
					Key:   key,
					Group: groups[key],
				}
				if !yield(group) {
					return
				}
			}
		},
	}
}
