package linq

func fromString[S ~string](source S) Query[rune] {
	return Query[rune]{iterate: func(yield func(rune) bool) {
		for _, value := range source {
			if !yield(value) {
				return
			}
		}
	}}
}
