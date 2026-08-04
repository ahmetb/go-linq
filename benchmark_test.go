package linq

import "testing"

const (
	size = 1000000
)

var benchmarkIntSlice []int

func BenchmarkTypedCore(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			var result []int
			legacyFromSlice(source).ToSlice(&result)
			benchmarkIntSlice = result
		}
	})

	b.Run("typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkIntSlice = fromSlice(source).toSlice()
		}
	})
}

func BenchmarkSelectWhereFirst(b *testing.B) {
	for n := 0; n < b.N; n++ {
		legacyRange(1, size).Select(func(i any) any {
			return -i.(int)
		}).Where(func(i any) bool {
			return i.(int) > -1000
		}).First()
	}
}

func BenchmarkSelectWhereFirst_generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		legacyRange(1, size).SelectT(func(i int) int {
			return -i
		}).WhereT(func(i int) bool {
			return i > -1000
		}).First()
	}
}

func BenchmarkSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		legacyRange(1, size).Where(func(i any) bool {
			return i.(int)%2 == 0
		}).SumInts()
	}
}

func BenchmarkSum_generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		legacyRange(1, size).WhereT(func(i int) bool {
			return i%2 == 0
		}).SumInts()
	}
}

func BenchmarkZipSkipTake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		legacyRange(1, size).Zip(legacyRange(1, size).Select(func(i any) any {
			return i.(int) * 2
		}), func(i, j any) any {
			return i.(int) + j.(int)
		}).Skip(2).Take(5)
	}
}

func BenchmarkZipSkipTake_generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		legacyRange(1, size).ZipT(legacyRange(1, size).SelectT(func(i int) int {
			return i * 2
		}), func(i, j int) int {
			return i + j
		}).Skip(2).Take(5)
	}
}
