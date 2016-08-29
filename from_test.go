package linq

import "testing"

func TestFrom(t *testing.T) {
	c := make(chan interface{}, 3)
	c <- -1
	c <- 0
	c <- 1
	close(c)

	tests := []struct {
		input  interface{}
		output []interface{}
		want   bool
	}{
		{[]int{1, 2, 3}, []interface{}{1, 2, 3}, true},
		{[]int{1, 2, 4}, []interface{}{1, 2, 3}, false},
		{[3]int{1, 2, 3}, []interface{}{1, 2, 3}, true},
		{[3]int{1, 2, 4}, []interface{}{1, 2, 3}, false},
		{"str", []interface{}{'s', 't', 'r'}, true},
		{"str", []interface{}{'s', 't', 'g'}, false},
		{map[string]bool{"foo": true}, []interface{}{KeyValue{"foo", true}}, true},
		{map[string]bool{"foo": true}, []interface{}{KeyValue{"foo", false}}, false},
		{c, []interface{}{-1, 0, 1}, true},
		{foo{f1: 1, f2: true, f3: "string"}, []interface{}{1, true, "string"}, true},
	}

	for _, test := range tests {
		if q := From(test.input); validateQuery(q, test.output) != test.want {
			if test.want {
				t.Errorf("From(%v)=%v expected %v", test.input, toSlice(q), test.output)
			} else {
				t.Errorf("From(%v)=%v expected not equal", test.input, test.output)
			}
		}
	}
}

func TestFromChannel(t *testing.T) {
	c := make(chan interface{}, 3)
	c <- 10
	c <- 15
	c <- -3
	close(c)

	w := []interface{}{10, 15, -3}

	if q := FromChannel(c); !validateQuery(q, w) {
		t.Errorf("FromChannel() failed expected %v", w)
	}
}

func TestFromString(t *testing.T) {
	s := "string"
	w := []interface{}{'s', 't', 'r', 'i', 'n', 'g'}

	if q := FromString(s); !validateQuery(q, w) {
		t.Errorf("FromString(%v)!=%v", s, w)
	}
}

func TestFromIterable(t *testing.T) {
	s := foo{f1: 1, f2: true, f3: "string"}
	w := []interface{}{1, true, "string"}

	if q := FromIterable(s); !validateQuery(q, w) {
		t.Errorf("FromIterable(%v)!=%v", s, w)
	}
}

func TestRange(t *testing.T) {
	w := []interface{}{-2, -1, 0, 1, 2}

	if q := Range(-2, 5); !validateQuery(q, w) {
		t.Errorf("Range(-2, 5)=%v expected %v", toSlice(q), w)
	}
}

func TestRepeat(t *testing.T) {
	w := []interface{}{1, 1, 1, 1, 1}

	if q := Repeat(1, 5); !validateQuery(q, w) {
		t.Errorf("Repeat(1, 5)=%v expected %v", toSlice(q), w)
	}
}
