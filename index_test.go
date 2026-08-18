package linq

import (
	"testing"
)

func TestIndexOf(t *testing.T) {
	arr := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	if index := FromSlice(arr[:]).IndexOf(func(i int) bool {
		return i == 3
	}); index != 2 {
		t.Errorf("IndexOf() expected 2 received %v", index)
	}

	if index := FromString("sstr").IndexOf(func(r rune) bool {
		return r == 'r'
	}); index != 3 {
		t.Errorf("IndexOf() expected 3 received %v", index)
	}

	if index := FromString("gadsgsadgsda").IndexOf(func(r rune) bool {
		return r == 'z'
	}); index != -1 {
		t.Errorf("IndexOf() expected -1 received %v", index)
	}
}

func TestElementAt(t *testing.T) {
	input := []int{10, 20, 30, 40, 50}
	tests := []struct {
		index  Index
		want   int
		wantOK bool
	}{
		{IndexFromStart(0), 10, true},
		{IndexFromStart(2), 30, true},
		{IndexFromStart(5), 0, false},
		{IndexFromEnd(1), 50, true},
		{IndexFromEnd(3), 30, true},
		{IndexFromEnd(5), 10, true},
		{IndexFromEnd(6), 0, false},
		{IndexFromEnd(0), 0, false},
	}

	for _, test := range tests {
		got, ok := FromSlice(input).ElementAt(test.index)
		if got != test.want || ok != test.wantOK {
			t.Errorf("ElementAt(%v)=%v,%v expected %v,%v",
				test.index, got, ok, test.want, test.wantOK)
		}
	}
}

func TestElementAtStopsAtRequestedIndex(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{10, 20, 30, 40}).Where(func(int) bool {
		pulled++
		return true
	})

	got, ok := q.ElementAt(IndexFromStart(2))
	if !ok || got != 30 {
		t.Fatalf("ElementAt(IndexFromStart(2))=%v,%v expected 30,true", got, ok)
	}
	if pulled != 3 {
		t.Errorf("source pulled %d elements expected 3", pulled)
	}
}

func TestIndexNegativeValuePanics(t *testing.T) {
	constructors := []func(){
		func() { NewIndex(-1, false) },
		func() { NewIndex(-1, true) },
		func() { IndexFromStart(-1) },
		func() { IndexFromEnd(-1) },
	}
	for _, constructor := range constructors {
		mustPanic(t, constructor)
	}
}
