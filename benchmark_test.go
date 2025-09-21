package linq

import "testing"

const (
	size = 1000000
)

func BenchmarkSelectWhereFirst(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).Select(func(i interface{}) interface{} {
			return -i.(int)
		}).Where(func(i interface{}) bool {
			return i.(int) > -1000
		}).First()
	}
}

func BenchmarkSelectWhereFirst_generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).SelectT(func(i int) int {
			return -i
		}).WhereT(func(i int) bool {
			return i > -1000
		}).First()
	}
}

func BenchmarkSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).Where(func(i interface{}) bool {
			return i.(int)%2 == 0
		}).SumInts()
	}
}

func BenchmarkSum_generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).WhereT(func(i int) bool {
			return i%2 == 0
		}).SumInts()
	}
}

func BenchmarkZipSkipTake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).Zip(Range(1, size).Select(func(i interface{}) interface{} {
			return i.(int) * 2
		}), func(i, j interface{}) interface{} {
			return i.(int) + j.(int)
		}).Skip(2).Take(5)
	}
}

func BenchmarkZipSkipTake_generics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).ZipT(Range(1, size).SelectT(func(i int) int {
			return i * 2
		}), func(i, j int) int {
			return i + j
		}).Skip(2).Take(5)
	}
}
