package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []interface{}{1, 2, 3, 4, 5}

	if q := From(input).Append(5); !validateQuery(q, want) {
		t.Errorf("From(%v).Append()=%v expected %v", input, toSlice(q), want)
	}
}

func TestAppendG(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{1, 2, 3, 4, 5}
	actual := FromSliceG(input).Append(5).ToSlice()
	assert.Equal(t, expected, actual)
}

func TestConcat(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{4, 5}
	want := []interface{}{1, 2, 3, 4, 5}

	if q := From(input1).Concat(From(input2)); !validateQuery(q, want) {
		t.Errorf("From(%v).Concat(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestConcatG(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{4, 5}
	expected := []int{1, 2, 3, 4, 5}
	actual := FromSliceG(input1).Concat(FromSliceG(input2)).ToSlice()
	assert.Equal(t, expected, actual)
}

func TestPrepend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []interface{}{0, 1, 2, 3, 4}

	if q := From(input).Prepend(0); !validateQuery(q, want) {
		t.Errorf("From(%v).Prepend()=%v expected %v", input, toSlice(q), want)
	}
}

func TestPrependG(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []int{0, 1, 2, 3, 4}
	actual := FromSliceG(input).Prepend(0).ToSlice()
	assert.Equal(t, want, actual)
}
