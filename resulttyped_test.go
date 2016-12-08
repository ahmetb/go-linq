package linq

import (
	"reflect"
	"testing"
)

func TestAllT(t *testing.T) {
	input := []int{2, 4, 6, 8}

	r1 := From(input).AllT(func(i int) bool {
		return i%2 == 0
	})
	r2 := From(input).AllT(func(i int) bool {
		return i%2 != 0
	})

	if !r1 {
		t.Errorf("From(%v).AllT()=%v", input, r1)
	}

	if r2 {
		t.Errorf("From(%v).AllT()=%v", input, r2)
	}
}

func TestAllT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AllT(func(item int) int { return item + 2 })
}

func TestAnyWithT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		want      bool
	}{
		{[]int{1, 2, 2, 3, 1}, func(i int) bool { return i == 4 }, false},
		{[]int{1, 2, 2, 3, 1}, func(i int) bool { return i == 4 }, false},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int) bool { return i == 4 }, true},
		{[]int{}, func(i int) bool { return i == 4 }, false},
	}

	for _, test := range tests {
		if r := From(test.input).AnyWithT(test.predicate); r != test.want {
			t.Errorf("From(%v).AnyWithT()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestAnyWithT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AnyWithT(func(item int) int { return item + 2 })
}

func TestCountWithT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		want      int
	}{
		{[]int{1, 2, 2, 3, 1}, func(i interface{}) bool { return i.(int) <= 2 }, 4},
		{[]int{1, 2, 2, 3, 1}, func(i int) bool { return i <= 2 }, 4},
		{[]int{}, func(i int) bool { return i <= 2 }, 0},
	}

	for _, test := range tests {
		if r := From(test.input).CountWithT(test.predicate); r != test.want {
			t.Errorf("From(%v).CountWithT()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestCountWithT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).CountWithT(func(item int) int { return item + 2 })
}

func TestFirstWithT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		want      interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, func(i interface{}) bool { return i.(int) > 2 }, 3},
		{[]int{1, 2, 2, 3, 1}, func(i int) bool { return i > 2 }, 3},
		{[]int{}, func(i interface{}) bool { return i.(int) > 2 }, nil},
	}

	for _, test := range tests {
		if r := From(test.input).FirstWithT(test.predicate); r != test.want {
			t.Errorf("From(%v).FirstWithT()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestFirstWithT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).FirstWithT(func(item int) int { return item + 2 })
}

func TestLastWithT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		want      interface{}
	}{
		{[]int{1, 2, 2, 3, 1, 4, 2, 5, 1, 1}, func(i interface{}) bool { return i.(int) > 2 }, 5},
		{[]int{1, 2, 2, 3, 1, 4, 2, 5, 1, 1}, func(i int) bool { return i > 2 }, 5},
		{[]int{}, func(i interface{}) bool { return i.(int) > 2 }, nil},
	}

	for _, test := range tests {
		if r := From(test.input).LastWithT(test.predicate); r != test.want {
			t.Errorf("From(%v).LastWith()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestLastWithT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).LastWithT(func(item int) int { return item + 2 })
}

func TestSingleWithT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		want      interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, func(i interface{}) bool { return i.(int) > 2 }, 3},
		{[]int{1, 1, 1}, func(i int) bool { return i > 2 }, nil},
		{[]int{5, 1, 1, 10, 2, 2}, func(i interface{}) bool { return i.(int) > 2 }, nil},
		{[]int{}, func(i interface{}) bool { return i.(int) > 2 }, nil},
	}

	for _, test := range tests {
		if r := From(test.input).SingleWithT(test.predicate); r != test.want {
			t.Errorf("From(%v).SingleWithT()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestSingleWithT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SingleWithT(func(item int) int { return item + 2 })
}

func TestToMapByT(t *testing.T) {
	input := make(map[int]bool)
	input[1] = true
	input[2] = false
	input[3] = true

	result := make(map[int]bool)
	From(input).ToMapByT(&result,
		func(i KeyValue) interface{} {
			return i.Key
		},
		func(i KeyValue) interface{} {
			return i.Value
		})

	if !reflect.DeepEqual(result, input) {
		t.Errorf("From(%v).ToMapByT()=%v expected %v", input, result, input)
	}
}

func TestToMapByT_PanicWhenKeySelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()
	result := make(map[int]bool)
	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).ToMapByT(
		&result,
		func(item, j int) int { return item + 2 },
		func(item int) int { return item + 2 },
	)
}

func TestToMapByT_PanicWhenValueSelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()
	result := make(map[int]bool)
	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).ToMapByT(
		&result,
		func(item int) int { return item + 2 },
		func(item, j int) int { return item + 2 },
	)
}
