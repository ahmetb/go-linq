package linq

import (
	"cmp"
	"crypto/sha256"
	"encoding/binary"
	"testing"
)

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

func BenchmarkChunk(b *testing.B) {
	const chunkSize = 256
	source := make([]int, 65536)
	b.Run("Count", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := Chunk(FromSlice(source), chunkSize).Count(); got != len(source)/chunkSize {
				b.Fatalf("Chunk count=%d expected %d", got, len(source)/chunkSize)
			}
		}
	})
	// Source setup is excluded; operator construction and consumption are timed.
	benchmark := func(b *testing.B, count, chunkSize int, unknown bool) {
		source := make([]int, count)
		for i := range source {
			source[i] = i + 1
		}
		q := FromSlice(source)
		if unknown {
			q = FromSeq(q.Iterate)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			total := 0
			Chunk(q, chunkSize).Iterate(func(chunk []int) bool {
				for _, v := range chunk {
					total += v
				}
				performanceChunkSink = chunk
				return true
			})
			performanceIntSink = total
		}
		b.StopTimer()
		if want := count * (count + 1) / 2; performanceIntSink != want {
			b.Fatalf("sum=%d want %d", performanceIntSink, want)
		}
	}
	b.Run("OversizedKnown", func(b *testing.B) { benchmark(b, 1, 1<<20, false) })
	b.Run("OversizedUnknown", func(b *testing.B) { benchmark(b, 1, 1<<20, true) })
	b.Run("FullKnown", func(b *testing.B) { benchmark(b, 65536, 256, false) })
	b.Run("FullUnknown", func(b *testing.B) { benchmark(b, 65536, 256, true) })
	b.Run("PartialKnown", func(b *testing.B) { benchmark(b, 65537, 256, false) })
}

func BenchmarkCountBy(b *testing.B) {
	const keys = 1024
	source := make([]int, 65536)
	for i := range source {
		source[i] = i % keys
	}
	q := FromSlice(source)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if got := q.CountBy(func(i int) int { return i }).Count(); got != keys {
			b.Fatalf("CountBy count=%d expected %d", got, keys)
		}
	}
}

func BenchmarkAggregateBy(b *testing.B) {
	const keys = 1024
	source := make([]int, 65536)
	for i := range source {
		source[i] = i % keys
	}
	q := FromSlice(source)
	keySelector := func(i int) int { return i }
	seedSelector := func(int) int { return 0 }
	accumulate := func(total, item int) int { return total + item }

	b.Run("Seed", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := q.AggregateBy(keySelector, 0, accumulate).Count(); got != keys {
				b.Fatalf("AggregateBy count=%d expected %d", got, keys)
			}
		}
	})
	b.Run("SeedSelector", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := q.AggregateByWithSeedSelector(keySelector, seedSelector, accumulate).Count(); got != keys {
				b.Fatalf("AggregateByWithSeedSelector count=%d expected %d", got, keys)
			}
		}
	})
}

func BenchmarkJoin(b *testing.B) {
	const keys = 1024
	source := make([]int, 65536)
	for i := range source {
		source[i] = i % keys
	}
	outer, inner := FromSlice(source), Range(0, keys)
	identity := func(i int) int { return i }
	pair := func(a, b int) int { return a + b }

	b.Run("Join", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := outer.Join(inner, identity, identity, pair).Count(); got != len(source) {
				b.Fatalf("Join count=%d expected %d", got, len(source))
			}
		}
	})
	b.Run("LeftJoin", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := outer.LeftJoin(inner, identity, identity, pair).Count(); got != len(source) {
				b.Fatalf("LeftJoin count=%d expected %d", got, len(source))
			}
		}
	})
	b.Run("RightJoin", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := inner.RightJoin(outer, identity, identity, pair).Count(); got != len(source) {
				b.Fatalf("RightJoin count=%d expected %d", got, len(source))
			}
		}
	})
	b.Run("FullJoin.Count", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := outer.FullJoin(inner, identity, identity, pair).Count(); got != len(source) {
				b.Fatalf("FullJoin count=%d expected %d", got, len(source))
			}
		}
	})
	benchmarkFullJoin := func(b *testing.B, outer, inner Query[int], first, matched bool) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q := outer.FullJoin(inner, identity, identity, pair)
			if first {
				q = q.Take(1)
			}
			performanceIntSink = Sum(q)
		}
		b.StopTimer()
		want := 65536 * 65537 / 2
		if first {
			want = 1
		}
		if matched {
			want *= 2
		}
		if performanceIntSink != want {
			b.Fatalf("sum=%d want %d", performanceIntSink, want)
		}
	}
	full := Range(1, 65536)
	emptyOuter := Empty[int]()
	emptyInner := FromSlice([]int{})
	b.Run("FullJoin.EmptyOuter", func(b *testing.B) { benchmarkFullJoin(b, emptyOuter, full, false, false) })
	b.Run("FullJoin.EmptyOuterFirst", func(b *testing.B) { benchmarkFullJoin(b, emptyOuter, full, true, false) })
	b.Run("FullJoin.EmptyInner", func(b *testing.B) { benchmarkFullJoin(b, full, emptyInner, false, false) })
	b.Run("FullJoin.EmptyInnerFirst", func(b *testing.B) { benchmarkFullJoin(b, full, emptyInner, true, false) })
	b.Run("FullJoin.Matched", func(b *testing.B) { benchmarkFullJoin(b, full, full, false, true) })
}

func BenchmarkSequence(b *testing.B) {
	const count = 65536
	q := Sequence(0, count-1, 1)

	for n := 0; n < b.N; n++ {
		if got := Sum(q); got != count*(count-1)/2 {
			b.Fatalf("Sequence sum=%d expected %d", got, count*(count-1)/2)
		}
	}
}

var performanceIntSink int
var performanceChunkSink []int

func BenchmarkChunkBy(b *testing.B) {
	benchmark := func(b *testing.B, count, size int, unknown bool) {
		source := make([]int, count)
		for i := range source {
			source[i] = i + 1
		}
		q := FromSlice(source)
		if unknown {
			q = FromSeq(q.Iterate)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			total := 0
			q.ChunkBy(size, func(chunk []int) int {
				part := 0
				for _, v := range chunk {
					part += v
				}
				return part
			}).Iterate(func(part int) bool { total += part; return true })
			performanceIntSink = total
		}
		b.StopTimer()
		if want := count * (count + 1) / 2; performanceIntSink != want {
			b.Fatalf("sum=%d want %d", performanceIntSink, want)
		}
	}
	b.Run("OversizedUnknown", func(b *testing.B) { benchmark(b, 1, 1<<20, true) })
	b.Run("FullKnown", func(b *testing.B) { benchmark(b, 65536, 256, false) })
}

func performanceOrderKey(v int) uint64 {
	var data [8]byte
	binary.LittleEndian.PutUint64(data[:], uint64(v))
	hash := sha256.Sum256(data[:])
	return binary.LittleEndian.Uint64(hash[:])
}

func benchmarkOrderingSource() Query[int] {
	source := make([]int, 4096)
	for i := range source {
		source[i] = (i * 7919) % len(source)
	}
	return FromSlice(source)
}

func benchmarkOrdering(b *testing.B, order func() Query[int]) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total := 0
		order().Iterate(func(v int) bool { total = total*31 + v; return true })
		performanceIntSink = total
	}
}

func BenchmarkOrderBy(b *testing.B) {
	q := benchmarkOrderingSource()
	b.Run("CheapKey", func(b *testing.B) {
		benchmarkOrdering(b, func() Query[int] { return q.OrderBy(func(v int) int { return v }).Query })
	})
	b.Run("ExpensiveKey", func(b *testing.B) {
		benchmarkOrdering(b, func() Query[int] { return q.OrderBy(performanceOrderKey).Query })
	})
}

func BenchmarkThenBy(b *testing.B) {
	q := benchmarkOrderingSource()
	benchmarkOrdering(b, func() Query[int] {
		return q.OrderBy(func(v int) int { return v % 8 }).ThenByDescending(performanceOrderKey).Query
	})
}

func BenchmarkOrderWith(b *testing.B) {
	q := benchmarkOrderingSource()
	benchmarkOrdering(b, func() Query[int] { return q.OrderWith(cmp.Compare[int]).Query })
}
