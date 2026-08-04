package linq

import (
	"slices"
	"testing"
)

func TestUnion(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []int{1, 2, 3, 4, 5}

	if q := FromSlice(input1).Union(FromSlice(input2)); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Union(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestUnionBy(t *testing.T) {
	type user struct {
		id   int
		name string
	}

	input1 := []user{{1, "Foo"}, {2, "Bar"}}
	input2 := []user{{3, "Foo"}, {4, "Baz"}}
	want := []user{{1, "Foo"}, {2, "Bar"}, {4, "Baz"}}

	if q := FromSlice(input1).UnionBy(FromSlice(input2), func(u user) string {
		return u.name
	}); !testQueryIteration(q, want) {
		t.Errorf("UnionBy()=%v expected %v", toSlice(q), want)
	}
}

func TestUnion_Abort(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}

	tests := []struct {
		name       string
		abortIndex int   // stop after this many items
		want       []int // expected collected values
	}{
		{
			name:       "iteration stops on input1",
			abortIndex: 2, // stops after 2 elements from input1
			want:       []int{1, 2},
		},
		{
			name:       "iteration stops on input2",
			abortIndex: 4, // stops after seeing the first new element from input2
			want:       []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := FromSlice(input1).Union(FromSlice(input2))

			var results []int
			i := 0
			q.Iterate(func(v int) bool {
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
