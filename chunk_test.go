package linq

import (
	"math"
	"slices"
	"testing"
)

func TestChunkOversizedAndUnknownSource(t *testing.T) {
	for _, source := range []Query[int]{
		FromSlice([]int{1, 2, 3}),
		FromSeq(FromSlice([]int{1, 2, 3}).Iterate),
	} {
		want := [][]int{{1, 2, 3}}
		if q := Chunk(source, math.MaxInt); !testQueryIteration(q, want) {
			t.Errorf("Chunk oversized result differs")
		}
		q := source.ChunkBy(math.MaxInt, func(chunk []int) []int { return chunk })
		if !testQueryIteration(q, want) {
			t.Errorf("ChunkBy oversized result differs")
		}
	}
}

func TestChunkGrowsPastStaleHint(t *testing.T) {
	source := map[int]int{1: 1}
	q := FromMap(source).Select(func(kv KeyValue[int, int]) int { return kv.Value })
	for i := 2; i <= 37; i++ {
		source[i] = i
	}
	chunks := Chunk(q, 9).ToSlice()
	var values []int
	for _, chunk := range chunks {
		if len(chunk) > 9 || cap(chunk) != len(chunk) {
			t.Errorf("invalid chunk len=%d cap=%d", len(chunk), cap(chunk))
		}
		values = append(values, chunk...)
	}
	slices.Sort(values)
	if !slices.Equal(values, Range(1, 37).ToSlice()) {
		t.Errorf("lost elements: %v", values)
	}
}

func TestChunkByRetainedSlices(t *testing.T) {
	chunks := FromSeq(Range(1, 19).Iterate).ChunkBy(9, func(chunk []int) []int { return chunk }).ToSlice()
	for i, chunk := range chunks {
		if cap(chunk) != len(chunk) {
			t.Errorf("chunk %d len=%d cap=%d", i, len(chunk), cap(chunk))
		}
		for j, v := range chunk {
			if v != i*9+j+1 {
				t.Errorf("chunk %d overwritten: %v", i, chunk)
			}
		}
	}
	if len(chunks) < 3 {
		t.Fatalf("got %d chunks expected 3", len(chunks))
	}
	chunks[0][0] = -1
	if chunks[1][0] != 10 || chunks[2][0] != 19 {
		t.Errorf("chunks share storage")
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		input []int
		size  int
		want  [][]int
	}{
		{nil, 3, nil},
		{[]int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
		{[]int{1, 2, 3, 4, 5, 6}, 3, [][]int{{1, 2, 3}, {4, 5, 6}}},
		{[]int{1, 2, 3, 4, 5, 6, 7}, 3, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}},
	}

	for _, test := range tests {
		if q := Chunk(FromSlice(test.input), test.size); !testQueryIteration(q, test.want) {
			t.Errorf("Chunk(%v, %d)=%v expected %v", test.input, test.size, toSlice(q), test.want)
		}
	}
}

func TestChunkClipsAndSeparatesSlices(t *testing.T) {
	chunks := Chunk(FromSlice([]int{1, 2, 3, 4, 5}), 3).ToSlice()
	for _, chunk := range chunks {
		if cap(chunk) != len(chunk) {
			t.Errorf("chunk %v has cap=%d expected %d", chunk, cap(chunk), len(chunk))
		}
	}

	chunks[0][0] = 99
	if !slices.Equal(chunks[1], []int{4, 5}) {
		t.Errorf("mutating one chunk changed another: %v", chunks)
	}
}

func TestChunkInvalidSizePanics(t *testing.T) {
	for _, size := range []int{0, -1} {
		mustPanic(t, func() { Chunk(FromSlice([]int{1}), size) })
	}
}

func TestChunkStopsPullingAfterConsumerStops(t *testing.T) {
	pulled := 0
	q := Chunk(FromSlice([]int{1, 2, 3, 4, 5}).Where(func(int) bool {
		pulled++
		return true
	}), 3)

	var got [][]int
	q.Iterate(func(chunk []int) bool {
		got = append(got, chunk)
		return false
	})

	if len(got) != 1 {
		t.Fatalf("Chunk early exit yielded %v expected one chunk", got)
	}
	if !slices.Equal(got[0], []int{1, 2, 3}) {
		t.Errorf("first chunk=%v expected [1 2 3]", got[0])
	}
	if pulled != 3 {
		t.Errorf("source pulled %d elements expected 3", pulled)
	}
}

func TestChunkBy(t *testing.T) {
	tests := []struct {
		input []int
		size  int
		want  []int
	}{
		{nil, 3, nil},
		{[]int{1, 2, 3}, 1, []int{1, 2, 3}},
		{[]int{1, 2, 3, 4, 5, 6}, 3, []int{6, 15}},
		{[]int{1, 2, 3, 4, 5, 6, 7}, 3, []int{6, 15, 7}},
	}

	sum := func(chunk []int) int {
		total := 0
		for _, v := range chunk {
			total += v
		}
		return total
	}

	for _, test := range tests {
		q := FromSlice(test.input).ChunkBy(test.size, sum)
		if !testQueryIteration(q, test.want) {
			t.Errorf("ChunkBy(%v, %d)=%v expected %v", test.input, test.size, toSlice(q), test.want)
		}
	}
}

func TestChunkByChainsAndChangesType(t *testing.T) {
	got := FromSlice([]int{1, 2, 3, 4, 5}).
		ChunkBy(2, func(chunk []int) int { return len(chunk) }).
		Where(func(n int) bool { return n == 2 }).
		ToSlice()

	if want := []int{2, 2}; !slices.Equal(got, want) {
		t.Errorf("ChunkBy().Where()=%v expected %v", got, want)
	}
}

func TestChunkByInvalidSizePanics(t *testing.T) {
	identity := func(chunk []int) []int { return chunk }
	for _, size := range []int{0, -1} {
		mustPanic(t, func() { FromSlice([]int{1}).ChunkBy(size, identity) })
	}
}

func TestChunkByStopsPullingAfterConsumerStops(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{1, 2, 3, 4, 5}).Where(func(int) bool {
		pulled++
		return true
	}).ChunkBy(3, func(chunk []int) int { return len(chunk) })

	var got []int
	q.Iterate(func(n int) bool {
		got = append(got, n)
		return false
	})

	if len(got) != 1 {
		t.Errorf("got %v expected one chunk", got)
		return
	}
	if pulled != 3 {
		t.Errorf("source pulled %d elements expected 3", pulled)
	}
}
