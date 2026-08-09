package linq

import "testing"

const (
	size = 1000000
)

func BenchmarkSelectWhereFirst(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).Select(func(i int) int {
			return -i
		}).Where(func(i int) bool {
			return i > -1000
		}).First()
	}
}

func BenchmarkSelectWhereFirst_handwritten(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var r int
		for i := 1; i <= size; i++ {
			v := -i
			if v > -1000 {
				r = v
				break
			}
		}
		_ = r
	}
}

func BenchmarkSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Sum(Range(1, size).Where(func(i int) bool {
			return i%2 == 0
		}))
	}
}

func BenchmarkSum_handwritten(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var sum int
		for i := 1; i <= size; i++ {
			if i%2 == 0 {
				sum += i
			}
		}
		_ = sum
	}
}

func BenchmarkZipSkipTake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).Zip(Range(1, size).Select(func(i int) int {
			return i * 2
		}), func(i, j int) int {
			return i + j
		}).Skip(2).Take(5).Count()
	}
}

// BenchmarkMaterializeSized covers the operators that compute a size hint
// from their inputs' counts: a regression in hint propagation shows up here
// as append-growth allocations instead of a single preallocation.
func BenchmarkMaterializeSized(b *testing.B) {
	src := make([]int, size)
	for i := range src {
		src[i] = i
	}
	tail := make([]int, size/10)

	b.Run("take", func(b *testing.B) {
		for b.Loop() {
			_ = FromSlice(src).Take(size / 2).ToSlice()
		}
	})
	b.Run("skip", func(b *testing.B) {
		for b.Loop() {
			_ = FromSlice(src).Skip(size / 2).ToSlice()
		}
	})
	b.Run("concat", func(b *testing.B) {
		for b.Loop() {
			_ = FromSlice(src).Concat(FromSlice(tail)).ToSlice()
		}
	})
	b.Run("prepend-append", func(b *testing.B) {
		for b.Loop() {
			_ = FromSlice(src).Append(1).Prepend(2).ToSlice()
		}
	})
	b.Run("zip", func(b *testing.B) {
		for b.Loop() {
			_ = FromSlice(src).Zip(FromSlice(src), func(a, z int) int { return a + z }).ToSlice()
		}
	})
}

// BenchmarkBuildMaps covers the map-building terminals and operators that
// presize their maps from a size hint.
func BenchmarkBuildMaps(b *testing.B) {
	src := make([]int, size)
	for i := range src {
		src[i] = i
	}
	inner := make([]int, size/100)
	for i := range inner {
		inner[i] = i
	}

	b.Run("tomapby", func(b *testing.B) {
		for b.Loop() {
			_ = FromSlice(src).ToMapBy(
				func(x int) int { return x },
				func(x int) int { return x * 2 })
		}
	})
	b.Run("join", func(b *testing.B) {
		for b.Loop() {
			_ = FromSlice(src[:size/10]).Join(FromSlice(inner),
				func(o int) int { return o % (size / 100) },
				func(i int) int { return i },
				func(o, i int) int { return o + i },
			).Count()
		}
	})
}
