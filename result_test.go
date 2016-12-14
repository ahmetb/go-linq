package linq

import (
	"math"
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	input := []int{2, 4, 6, 8}

	r1 := From(input).All(func(i interface{}) bool {
		return i.(int)%2 == 0
	})
	r2 := From(input).All(func(i interface{}) bool {
		return i.(int)%2 != 0
	})

	if !r1 {
		t.Errorf("From(%v).All()=%v", input, r1)
	}

	if r2 {
		t.Errorf("From(%v).All()=%v", input, r2)
	}
}

func TestAllT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "AllT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AllT(func(item int) int { return item + 2 })
	})
}

func TestAny(t *testing.T) {
	tests := []struct {
		input interface{}
		want  bool
	}{
		{[]int{1, 2, 2, 3, 1}, true},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, true},
		{"sstr", true},
		{[]int{}, false},
	}

	for _, test := range tests {
		if r := From(test.input).Any(); r != test.want {
			t.Errorf("From(%v).Any()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestAnyWith(t *testing.T) {
	tests := []struct {
		input interface{}
		want  bool
	}{
		{[]int{1, 2, 2, 3, 1}, false},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, true},
		{[]int{}, false},
	}

	for _, test := range tests {
		if r := From(test.input).AnyWith(func(i interface{}) bool {
			return i.(int) == 4
		}); r != test.want {
			t.Errorf("From(%v).Any()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestAnyWithT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "AnyWithT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AnyWithT(func(item int) int { return item + 2 })
	})
}

func TestAverage(t *testing.T) {
	tests := []struct {
		input interface{}
		want  float64
	}{
		{[]int{1, 2, 2, 3, 1}, 1.8},
		{[5]uint{1, 2, 5, 7, 10}, 5.},
		{[]float32{1., 1.}, 1.},
	}

	for _, test := range tests {
		if r := From(test.input).Average(); r != test.want {
			t.Errorf("From(%v).Average()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestAverageForNaN(t *testing.T) {
	if r := From([]int{}).Average(); !math.IsNaN(r) {
		t.Errorf("From([]int{}).Average()=%v expected %v", r, math.NaN())
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		input interface{}
		value interface{}
		want  bool
	}{
		{[]int{1, 2, 2, 3, 1}, 10, false},
		{[5]uint{1, 2, 5, 7, 10}, uint(5), true},
		{[]float32{}, 1., false},
	}

	for _, test := range tests {
		if r := From(test.input).Contains(test.value); r != test.want {
			t.Errorf("From(%v).Contains(%v)=%v expected %v", test.input, test.value, r, test.want)
		}
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		input interface{}
		want  int
	}{
		{[]int{1, 2, 2, 3, 1}, 5},
		{[7]uint{1, 2, 5, 7, 10, 12, 15}, 7},
		{[]float32{}, 0},
	}

	for _, test := range tests {
		if r := From(test.input).Count(); r != test.want {
			t.Errorf("From(%v).Count()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestCountWith(t *testing.T) {
	tests := []struct {
		input interface{}
		want  int
	}{
		{[]int{1, 2, 2, 3, 1}, 4},
		{[]int{}, 0},
	}

	for _, test := range tests {
		if r := From(test.input).CountWith(func(i interface{}) bool {
			return i.(int) <= 2
		}); r != test.want {
			t.Errorf("From(%v).CountWith()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestCountWithT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "CountWithT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).CountWithT(func(item int) int { return item + 2 })
	})
}

func TestFirst(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, 1},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).First(); r != test.want {
			t.Errorf("From(%v).First()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestFirstWith(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, 3},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).FirstWith(func(i interface{}) bool {
			return i.(int) > 2
		}); r != test.want {
			t.Errorf("From(%v).FirstWith()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestFirstWithT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "FirstWithT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).FirstWithT(func(item int) int { return item + 2 })
	})
}

func TestLast(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, 1},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).Last(); r != test.want {
			t.Errorf("From(%v).Last()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestLastWith(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 1, 4, 2, 5, 1, 1}, 5},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).LastWith(func(i interface{}) bool {
			return i.(int) > 2
		}); r != test.want {
			t.Errorf("From(%v).LastWith()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestLastWithT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "LastWithT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).LastWithT(func(item int) int { return item + 2 })
	})
}

func TestMax(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, 3},
		{[]int{1}, 1},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).Max(); r != test.want {
			t.Errorf("From(%v).Max()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 0}, 0},
		{[]int{1}, 1},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).Min(); r != test.want {
			t.Errorf("From(%v).Min()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestResults(t *testing.T) {
	input := []int{1, 2, 3}
	want := []interface{}{1, 2, 3}

	if r := From(input).Results(); !reflect.DeepEqual(r, want) {
		t.Errorf("From(%v).Raw()=%v expected %v", input, r, want)
	}
}

func TestSequenceEqual(t *testing.T) {
	tests := []struct {
		input  interface{}
		input2 interface{}
		want   bool
	}{
		{[]int{1, 2, 2, 3, 1}, []int{4, 6}, false},
		{[]int{1, -1, 100}, []int{1, -1, 100}, true},
		{[]int{}, []int{}, true},
	}

	for _, test := range tests {
		if r := From(test.input).SequenceEqual(From(test.input2)); r != test.want {
			t.Errorf("From(%v).SequenceEqual(%v)=%v expected %v", test.input, test.input2, r, test.want)
		}
	}
}

func TestSingle(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, nil},
		{[]int{1}, 1},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).Single(); r != test.want {
			t.Errorf("From(%v).Single()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestSingleWith(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, 3},
		{[]int{1, 1, 1}, nil},
		{[]int{5, 1, 1, 10, 2, 2}, nil},
		{[]int{}, nil},
	}

	for _, test := range tests {
		if r := From(test.input).SingleWith(func(i interface{}) bool {
			return i.(int) > 2
		}); r != test.want {
			t.Errorf("From(%v).SingleWith()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestSingleWithT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SingleWithT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SingleWithT(func(item int) int { return item + 2 })
	})
}

func TestSumInts(t *testing.T) {
	tests := []struct {
		input interface{}
		want  int64
	}{
		{[]int{1, 2, 2, 3, 1}, 9},
		{[]int{1}, 1},
		{[]int{}, 0},
	}

	for _, test := range tests {
		if r := From(test.input).SumInts(); r != test.want {
			t.Errorf("From(%v).SumInts()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestSumUInts(t *testing.T) {
	tests := []struct {
		input interface{}
		want  uint64
	}{
		{[]uint{1, 2, 2, 3, 1}, 9},
		{[]uint{1}, 1},
		{[]uint{}, 0},
	}

	for _, test := range tests {
		if r := From(test.input).SumUInts(); r != test.want {
			t.Errorf("From(%v).SumInts()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestSumFloats(t *testing.T) {
	tests := []struct {
		input interface{}
		want  float64
	}{
		{[]float32{1., 2., 2., 3., 1.}, 9.},
		{[]float64{1.}, 1.},
		{[]float32{}, 0.},
	}

	for _, test := range tests {
		if r := From(test.input).SumFloats(); r != test.want {
			t.Errorf("From(%v).SumFloats()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestToChannel(t *testing.T) {
	c := make(chan interface{})
	input := []int{1, 2, 3, 4, 5}

	go func() {
		From(input).ToChannel(c)
	}()

	result := []int{}
	for value := range c {
		result = append(result, value.(int))
	}

	if !reflect.DeepEqual(result, input) {
		t.Errorf("From(%v).ToChannel()=%v expected %v", input, result, input)
	}
}

func TestToMap(t *testing.T) {
	input := make(map[int]bool)
	input[1] = true
	input[2] = false
	input[3] = true

	result := make(map[int]bool)
	From(input).ToMap(&result)

	if !reflect.DeepEqual(result, input) {
		t.Errorf("From(%v).ToMap()=%v expected %v", input, result, input)
	}
}

func TestToMapBy(t *testing.T) {
	input := make(map[int]bool)
	input[1] = true
	input[2] = false
	input[3] = true

	result := make(map[int]bool)
	From(input).ToMapBy(&result,
		func(i interface{}) interface{} {
			return i.(KeyValue).Key
		},
		func(i interface{}) interface{} {
			return i.(KeyValue).Value
		})

	if !reflect.DeepEqual(result, input) {
		t.Errorf("From(%v).ToMapBy()=%v expected %v", input, result, input)
	}
}

func TestToMapByT_PanicWhenKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "ToMapByT: parameter [keySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		result := make(map[int]bool)
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).ToMapByT(
			&result,
			func(item, j int) int { return item + 2 },
			func(item int) int { return item + 2 },
		)
	})
}

func TestToMapByT_PanicWhenValueSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "ToMapByT: parameter [valueSelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		result := make(map[int]bool)
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).ToMapByT(
			&result,
			func(item int) int { return item + 2 },
			func(item, j int) int { return item + 2 },
		)
	})
}

func TestToSlice(t *testing.T) {
	input := []int{1, 2, 3, 4}

	result := []int{}
	From(input).ToSlice(&result)

	if !reflect.DeepEqual(result, input) {
		t.Errorf("From(%v).ToSlice()=%v expected %v", input, result, input)
	}
}
