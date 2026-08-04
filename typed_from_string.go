package linq

func fromString[S ~string](source S) query[rune] {
	return query[rune]{iterate: func(yield func(rune) bool) {
		for _, value := range source {
			if !yield(value) {
				return
			}
		}
	}}
}
