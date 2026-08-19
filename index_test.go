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
		index  Position
		want   int
		wantOK bool
	}{
		{PositionFromStart(0), 10, true},
		{PositionFromStart(2), 30, true},
		{PositionFromStart(5), 0, false},
		{PositionFromEnd(1), 50, true},
		{PositionFromEnd(3), 30, true},
		{PositionFromEnd(5), 10, true},
		{PositionFromEnd(6), 0, false},
		{PositionFromEnd(0), 0, false},
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

	got, ok := q.ElementAt(PositionFromStart(2))
	if !ok || got != 30 {
		t.Fatalf("ElementAt(PositionFromStart(2))=%v,%v expected 30,true", got, ok)
	}
	if pulled != 3 {
		t.Errorf("source pulled %d elements expected 3", pulled)
	}
}

func TestPositionNegativeValuePanics(t *testing.T) {
	constructors := []func(){
		func() { NewPosition(-1, false) },
		func() { NewPosition(-1, true) },
		func() { PositionFromStart(-1) },
		func() { PositionFromEnd(-1) },
	}
	for _, constructor := range constructors {
		mustPanic(t, constructor)
	}
}

func TestIndex(t *testing.T) {
	want := []KeyValue[int, string]{
		{Key: 0, Value: "zero"},
		{Key: 1, Value: "one"},
		{Key: 2, Value: "two"},
	}

	if q := Index(FromSlice([]string{"zero", "one", "two"})); !testQueryIteration(q, want) {
		t.Errorf("Index()=%v expected %v", q.ToSlice(), want)
	}
	if got := Index(FromSlice([]string{})).ToSlice(); got != nil {
		t.Errorf("Index(empty)=%v expected nil", got)
	}
}

func TestIndexStopsWithConsumer(t *testing.T) {
	pulled := 0
	q := Index(FromSlice([]int{10, 20, 30}).Where(func(int) bool {
		pulled++
		return true
	}))

	q.Iterate(func(KeyValue[int, int]) bool { return false })
	if pulled != 1 {
		t.Errorf("Index pulled %d elements expected 1", pulled)
	}
}
