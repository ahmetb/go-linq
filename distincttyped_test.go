package linq

import "testing"

func TestDistinctByT(t *testing.T) {
	type user struct {
		id   int
		name string
	}

	tests := []struct {
		input    interface{}
		selector interface{}
		output   []interface{}
	}{
		{[]user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}, func(u interface{}) interface{} {
			return u.(user).name
		}, []interface{}{user{1, "Foo"}, user{2, "Bar"}}},
		{[]user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}, func(u user) string {
			return u.name
		}, []interface{}{user{1, "Foo"}, user{2, "Bar"}}},
	}

	for _, test := range tests {
		if q := From(test.input).DistinctByT(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).WhereIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestDistinctBy_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).DistinctByT(func(indice, item string) bool { return item == "2" })
}
