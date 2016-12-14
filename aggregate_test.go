package linq

import "testing"
import "strings"

func TestAggregate(t *testing.T) {
	tests := []struct {
		input interface{}
		want  interface{}
	}{
		{[]string{"apple", "mango", "orange", "passionfruit", "grape"}, "passionfruit"},
		{[]string{}, nil},
	}

	for _, test := range tests {
		r := From(test.input).Aggregate(func(r interface{}, i interface{}) interface{} {
			if len(r.(string)) > len(i.(string)) {
				return r
			}
			return i
		})

		if r != test.want {
			t.Errorf("From(%v).Aggregate()=%v expected %v", test.input, r, test.want)
		}
	}
}

func TestAggregateT_PanicWhenFunctionIsInvalid(t *testing.T) {
	mustPanicWithError(t, "AggregateT: parameter [f] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,string,string)string'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AggregateT(func(x int, r string, i string) string {
			if len(r) > len(i) {
				return r
			}
			return i
		})
	})
}

func TestAggregateWithSeed(t *testing.T) {
	input := []string{"apple", "mango", "orange", "banana", "grape"}
	want := "passionfruit"

	r := From(input).AggregateWithSeed(want,
		func(r interface{}, i interface{}) interface{} {
			if len(r.(string)) > len(i.(string)) {
				return r
			}
			return i
		})

	if r != want {
		t.Errorf("From(%v).AggregateWithSeed()=%v expected %v", input, r, want)
	}
}

func TestAggregateWithSeedT_PanicWhenFunctionIsInvalid(t *testing.T) {
	mustPanicWithError(t, "AggregateWithSeed: parameter [f] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,string,string)string'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AggregateWithSeedT(3, func(x int, r string, i string) string {
			if len(r) > len(i) {
				return r
			}
			return i
		})
	})
}

func TestAggregateWithSeedBy(t *testing.T) {
	input := []string{"apple", "mango", "orange", "passionfruit", "grape"}
	want := "PASSIONFRUIT"

	r := From(input).AggregateWithSeedBy("banana",
		func(r interface{}, i interface{}) interface{} {
			if len(r.(string)) > len(i.(string)) {
				return r
			}
			return i
		},
		func(r interface{}) interface{} {
			return strings.ToUpper(r.(string))
		},
	)

	if r != want {
		t.Errorf("From(%v).AggregateWithSeed()=%v expected %v", input, r, want)
	}
}

func TestAggregateWithSeedByT_PanicWhenFunctionIsInvalid(t *testing.T) {
	mustPanicWithError(t, "AggregateWithSeedByT: parameter [f] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,string,string)string'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AggregateWithSeedByT(3,
			func(x int, r string, i string) string {
				if len(r) > len(i) {
					return r
				}
				return i
			},
			func(r string) string {
				return r
			},
		)
	})
}

func TestAggregateWithSeedByT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "AggregateWithSeedByT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(string,int)string'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).AggregateWithSeedByT(3,
			func(x int, r int) int {
				if x > r {
					return x
				}
				return r
			},
			func(r string, t int) string {
				return r
			},
		)
	})
}
