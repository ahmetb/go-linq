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
