package linq

import "testing"

func TestZip(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []int{3, 6, 8}

	if q := FromSlice(input1).Zip(FromSlice(input2), func(i, j int) int {
		return i + j
	}); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Zip(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestZip_TypeChanging(t *testing.T) {
	input1 := []string{"a", "b", "c"}
	input2 := []int{1, 2, 3}
	want := []string{"a1", "b2", "c3"}

	if q := FromSlice(input1).Zip(FromSlice(input2), func(s string, i int) string {
		return s + string(rune('0'+i))
	}); !testQueryIteration(q, want) {
		t.Errorf("Zip()=%v expected %v", toSlice(q), want)
	}
}

func TestZip3(t *testing.T) {
	type triple struct {
		number int
		word   string
		flag   bool
	}

	tests := []struct {
		numbers []int
		words   []string
		flags   []bool
		want    []triple
	}{
		// first is shortest
		{[]int{1}, []string{"one", "two"}, []bool{true, false},
			[]triple{{1, "one", true}}},
		// second is shortest
		{[]int{1, 2, 3}, []string{"one"}, []bool{true, false, false},
			[]triple{{1, "one", true}}},
		// third is shortest
		{[]int{1, 2, 3}, []string{"one", "two", "three", "four"}, []bool{true, false},
			[]triple{{1, "one", true}, {2, "two", false}}},
	}

	for _, test := range tests {
		q := FromSlice(test.numbers).Zip3(
			FromSlice(test.words),
			FromSlice(test.flags),
			func(number int, word string, flag bool) triple {
				return triple{number, word, flag}
			},
		)
		if !testQueryIteration(q, test.want) {
			t.Errorf("Zip3(%v, %v, %v)=%v expected %v",
				test.numbers, test.words, test.flags, toSlice(q), test.want)
		}
	}
}

func TestZip3StopsAllSources(t *testing.T) {
	pulledSecond := 0
	pulledThird := 0
	second := FromSlice([]int{10, 20, 30}).Where(func(int) bool {
		pulledSecond++
		return true
	})
	third := FromSlice([]int{100, 200, 300}).Where(func(int) bool {
		pulledThird++
		return true
	})

	var got []int
	FromSlice([]int{1, 2, 3}).Zip3(second, third, func(a, b, c int) int {
		return a + b + c
	}).Iterate(func(item int) bool {
		got = append(got, item)
		return false
	})

	if len(got) != 1 || got[0] != 111 {
		t.Errorf("Zip3 early exit yielded %v expected [111]", got)
	}
	if pulledSecond != 1 || pulledThird != 1 {
		t.Errorf("Zip3 pulled second=%d third=%d expected 1 each", pulledSecond, pulledThird)
	}
}
