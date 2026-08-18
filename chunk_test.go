package linq

import (
	"slices"
	"testing"
)

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
