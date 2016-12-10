package linq

import "testing"

func TestZip(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []interface{}{3, 6, 8}

	if q := From(input1).Zip(From(input2), func(i, j interface{}) interface{} {
		return i.(int) + j.(int)
	}); !validateQuery(q, want) {
		t.Errorf("From(%v).Zip(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestZipT(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []interface{}{3, 6, 8}

	if q := From(input1).ZipT(From(input2), func(i, j interface{}) interface{} {
		return i.(int) + j.(int)
	}); !validateQuery(q, want) {
		t.Errorf("From(%v).Zip(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}

	if q := From(input1).ZipT(From(input2), func(i, j int) interface{} {
		return i + j
	}); !validateQuery(q, want) {
		t.Errorf("From(%v).Zip(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestZipT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}

	From(input1).ZipT(From(input2), func(i, j, k int) int {
		return i + j
	})

}
