package linq

import (
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	input := []int{2, 4, 6, 8}

	r1 := FromSlice(input).All(func(i int) bool {
		return i%2 == 0
	})
	r2 := FromSlice(input).All(func(i int) bool {
		return i%2 != 0
	})

	if !r1 {
		t.Errorf("FromSlice(%v).All()=%v", input, r1)
	}

	if r2 {
		t.Errorf("FromSlice(%v).All()=%v", input, r2)
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		input []int
		want  bool
	}{
		{[]int{1, 2, 2, 3, 1}, true},
		{[]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, true},
		{[]int{}, false},
	}

	for _, test := range tests {
		if r := FromSlice(test.input).Any(); r != test.want {
			t.Errorf("FromSlice(%v).Any()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestAnyWith(t *testing.T) {
	tests := []struct {
		input []int
		want  bool
	}{
		{[]int{1, 2, 2, 3, 1}, false},
		{[]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, true},
		{[]int{}, false},
	}

	for _, test := range tests {
		if r := FromSlice(test.input).AnyWith(func(i int) bool {
			return i == 4
		}); r != test.want {
			t.Errorf("FromSlice(%v).AnyWith()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestContains(t *testing.T) {
	if r := FromSlice([]int{1, 2, 2, 3, 1}).Contains(10); r {
		t.Errorf("Contains(10)=%v expected false", r)
	}

	if r := FromSlice([]uint{1, 2, 5, 7, 10}).Contains(5); !r {
		t.Errorf("Contains(5)=%v expected true", r)
	}

	if r := FromSlice([]float32{}).Contains(1.); r {
		t.Errorf("Contains(1.)=%v expected false", r)
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		input []int
		want  int
	}{
		{[]int{1, 2, 2, 3, 1}, 5},
		{[]int{1, 2, 5, 7, 10, 12, 15}, 7},
		{[]int{}, 0},
	}

	for _, test := range tests {
		if r := FromSlice(test.input).Count(); r != test.want {
			t.Errorf("FromSlice(%v).Count()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestCountWith(t *testing.T) {
	tests := []struct {
		input []int
		want  int
	}{
		{[]int{1, 2, 2, 3, 1}, 4},
		{[]int{}, 0},
	}

	for _, test := range tests {
		if r := FromSlice(test.input).CountWith(func(i int) bool {
			return i <= 2
		}); r != test.want {
			t.Errorf("FromSlice(%v).CountWith()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestFirst(t *testing.T) {
	if r, ok := FromSlice([]int{1, 2, 2, 3, 1}).First(); !ok || r != 1 {
		t.Errorf("First()=%v,%v expected 1,true", r, ok)
	}

	if r, ok := FromSlice([]int{}).First(); ok || r != 0 {
		t.Errorf("First()=%v,%v expected 0,false", r, ok)
	}
}

func TestFirstWith(t *testing.T) {
	if r, ok := FromSlice([]int{1, 2, 2, 3, 1}).FirstWith(func(i int) bool {
		return i > 2
	}); !ok || r != 3 {
		t.Errorf("FirstWith()=%v,%v expected 3,true", r, ok)
	}

	if r, ok := FromSlice([]int{}).FirstWith(func(i int) bool {
		return i > 2
	}); ok || r != 0 {
		t.Errorf("FirstWith()=%v,%v expected 0,false", r, ok)
	}
}

func TestForEach(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 2, 2, 35, 111}, []int{2, 4, 4, 70, 222}},
		{[]int{}, []int{}},
	}

	for _, test := range tests {
		output := []int{}
		FromSlice(test.input).ForEach(func(item int) {
			output = append(output, item*2)
		})

		if !reflect.DeepEqual(output, test.want) {
			t.Fatalf("FromSlice(%#v).ForEach()=%#v expected=%#v", test.input, output, test.want)
		}
	}
}

func TestForEachIndexed(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 2, 2, 35, 111}, []int{1, 3, 4, 38, 115}},
		{[]int{}, []int{}},
	}

	for _, test := range tests {
		output := []int{}
		FromSlice(test.input).ForEachIndexed(func(index int, item int) {
			output = append(output, item+index)
		})

		if !reflect.DeepEqual(output, test.want) {
			t.Fatalf("FromSlice(%#v).ForEachIndexed()=%#v expected=%#v", test.input, output, test.want)
		}
	}
}

func TestLast(t *testing.T) {
	if r, ok := FromSlice([]int{1, 2, 2, 3, 5}).Last(); !ok || r != 5 {
		t.Errorf("Last()=%v,%v expected 5,true", r, ok)
	}

	if r, ok := FromSlice([]int{}).Last(); ok || r != 0 {
		t.Errorf("Last()=%v,%v expected 0,false", r, ok)
	}
}

func TestLastWith(t *testing.T) {
	if r, ok := FromSlice([]int{1, 2, 2, 3, 1, 4, 2, 5, 1, 1}).LastWith(func(i int) bool {
		return i > 2
	}); !ok || r != 5 {
		t.Errorf("LastWith()=%v,%v expected 5,true", r, ok)
	}

	if r, ok := FromSlice([]int{}).LastWith(func(i int) bool {
		return i > 2
	}); ok || r != 0 {
		t.Errorf("LastWith()=%v,%v expected 0,false", r, ok)
	}
}

func TestSequenceEqual(t *testing.T) {
	tests := []struct {
		input  []int
		input2 []int
		want   bool
	}{
		{[]int{1, 2, 2, 3, 1}, []int{4, 6}, false},
		{[]int{1, -1, 100}, []int{1, -1, 100}, true},
		{[]int{}, []int{}, true},
		{[]int{1, 2}, []int{1, 2, 3}, false},
		{[]int{1, 2, 3}, []int{1, 2}, false},
	}

	for _, test := range tests {
		if r := FromSlice(test.input).SequenceEqual(FromSlice(test.input2)); r != test.want {
			t.Errorf("FromSlice(%v).SequenceEqual(%v)=%v expected %v", test.input, test.input2, r, test.want)
		}
	}
}

func TestSingle(t *testing.T) {
	tests := []struct {
		input  []int
		want   int
		wantOk bool
	}{
		{[]int{1, 2, 2, 3, 1}, 0, false},
		{[]int{1}, 1, true},
		{[]int{}, 0, false},
	}

	for _, test := range tests {
		if r, ok := FromSlice(test.input).Single(); ok != test.wantOk || r != test.want {
			t.Errorf("FromSlice(%v).Single()=%v,%v expected %v,%v", test.input, r, ok, test.want, test.wantOk)
		}
	}
}

func TestSingleWith(t *testing.T) {
	tests := []struct {
		input  []int
		want   int
		wantOk bool
	}{
		{[]int{1, 2, 2, 3, 1}, 3, true},
		{[]int{1, 1, 1}, 0, false},
		{[]int{5, 1, 1, 10, 2, 2}, 0, false},
		{[]int{}, 0, false},
	}

	for _, test := range tests {
		if r, ok := FromSlice(test.input).SingleWith(func(i int) bool {
			return i > 2
		}); ok != test.wantOk || r != test.want {
			t.Errorf("FromSlice(%v).SingleWith()=%v,%v expected %v,%v", test.input, r, ok, test.want, test.wantOk)
		}
	}
}

func TestToChannel(t *testing.T) {
	c := make(chan int)
	input := []int{1, 2, 3, 4, 5}

	go func() {
		FromSlice(input).ToChannel(c)
	}()

	result := []int{}
	for value := range c {
		result = append(result, value)
	}

	if !reflect.DeepEqual(result, input) {
		t.Errorf("FromSlice(%v).ToChannel()=%v expected %v", input, result, input)
	}
}

func TestToMap(t *testing.T) {
	input := make(map[int]bool)
	input[1] = true
	input[2] = false
	input[3] = true

	result := ToMap(FromMap(input))

	if !reflect.DeepEqual(result, input) {
		t.Errorf("ToMap(FromMap(%v))=%v expected %v", input, result, input)
	}
}

func TestToMapBy(t *testing.T) {
	input := make(map[int]bool)
	input[1] = true
	input[2] = false
	input[3] = true

	result := FromMap(input).ToMapBy(
		func(i KeyValue[int, bool]) int {
			return i.Key
		},
		func(i KeyValue[int, bool]) bool {
			return i.Value
		})

	if !reflect.DeepEqual(result, input) {
		t.Errorf("FromMap(%v).ToMapBy()=%v expected %v", input, result, input)
	}
}

func TestToMapBy_TypeChanging(t *testing.T) {
	input := []string{"apple", "banana"}
	want := map[string]int{"apple": 5, "banana": 6}

	result := FromSlice(input).ToMapBy(
		func(s string) string { return s },
		func(s string) int { return len(s) })

	if !reflect.DeepEqual(result, want) {
		t.Errorf("ToMapBy()=%v expected %v", result, want)
	}
}

func TestToSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := FromSlice(input).ToSlice()
	if !reflect.DeepEqual(result, input) {
		t.Errorf("FromSlice(%v).ToSlice()=%v expected %v", input, result, input)
	}

	var empty []int
	if result := FromSlice(empty).ToSlice(); len(result) != 0 {
		t.Errorf("FromSlice(nil).ToSlice()=%v expected empty", result)
	}
}

func TestToSlice_ReturnsNilWhenSourceMapIsCleared(t *testing.T) {
	source := map[int]int{1: 1}
	q := FromMap(source)
	clear(source)

	if out := q.ToSlice(); out != nil {
		t.Errorf("ToSlice() returned len=%d cap=%d; want nil", len(out), cap(out))
	}
}

func TestToSlice_ReturnsNilWhenEmptyWithoutSizeHint(t *testing.T) {
	// FromString carries no size hint, so this exercises the slices.Collect
	// path of collect().
	if out := FromString("").ToSlice(); out != nil {
		t.Errorf("ToSlice() returned len=%d cap=%d; want nil", len(out), cap(out))
	}
}
