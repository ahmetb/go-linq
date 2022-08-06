package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnion(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []interface{}{1, 2, 3, 4, 5}

	if q := From(input1).Union(From(input2)); !validateQuery(q, want) {
		t.Errorf("From(%v).Union(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestUnionG_int(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []int{1, 2, 3, 4, 5}

	assert.Equal(t, want, FromSliceG(input1).Union(FromSliceG(input2)).ToSlice())
}

type unionG_test struct {
	f1 int
	f2 string
}

func TestUnionG_struct(t *testing.T) {
	input1 := []unionG_test{
		{1, "1"},
		{2, "2"},
	}
	input2 := []unionG_test{
		{2, "2"},
		{3, "3"},
		{4, "4"},
		{5, "5"},
	}
	want := []unionG_test{
		{1, "1"},
		{2, "2"},
		{3, "3"},
		{4, "4"},
		{5, "5"},
	}

	assert.Equal(t, want, FromSliceG(input1).Union(FromSliceG(input2)).ToSlice())
}
