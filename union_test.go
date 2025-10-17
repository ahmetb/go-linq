package linq

import (
	"slices"
	"testing"
)

func TestUnion(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []any{1, 2, 3, 4, 5}

	if q := From(input1).Union(From(input2)); !testQueryIteration(q, want) {
		t.Errorf("From(%v).Union(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestUnion_Abort(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}

	tests := []struct {
		name       string
		abortIndex int   // stop after this many items
		want       []any // expected collected values
	}{
		{
			name:       "iteration stops on input1",
			abortIndex: 2, // stops after 2 elements from input1
			want:       []any{1, 2},
		},
		{
			name:       "iteration stops on input2",
			abortIndex: 4, // stops after seeing the first new element from input2
			want:       []any{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := From(input1).Union(From(input2))

			var results []any
			i := 0
			q.Iterate(func(v any) bool {
				results = append(results, v)
				i++
				if i >= tt.abortIndex {
					return false // simulate early termination
				}
				return true
			})

			if !slices.Equal(results, tt.want) {
				t.Errorf("%s: got %v, want %v", tt.name, results, tt.want)
			}
		})
	}
}
