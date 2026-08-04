package linq

func integerRange(start, count int) query[int] {
	return query[int]{iterate: func(yield func(int) bool) {
		for offset := 0; offset < count; offset++ {
			if !yield(start + offset) {
				return
			}
		}
	}}
}
