package linq

import "testing"

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
