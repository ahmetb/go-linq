package linq

import "testing"

func TestJoin(t *testing.T) {
	outer := []int{0, 1, 2, 3, 4, 5, 8}
	inner := []int{1, 2, 1, 4, 7, 6, 7, 2}
	want := []KeyValue[int, int]{
		{1, 1},
		{1, 1},
		{2, 2},
		{2, 2},
		{4, 4},
	}

	q := FromSlice(outer).Join(
		FromSlice(inner),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer int, inner int) KeyValue[int, int] {
			return KeyValue[int, int]{outer, inner}
		})

	if !testQueryIteration(q, want) {
		t.Errorf("FromSlice().Join()=%v expected %v", toSlice(q), want)
	}
}

func TestJoin_TypeChanging(t *testing.T) {
	type person struct {
		name   string
		petIDs []int
	}
	type pet struct {
		id   int
		name string
	}

	people := []person{{"Ahmet", []int{1, 2}}, {"Bora", []int{3}}}
	pets := []pet{{1, "Fluffy"}, {2, "Rex"}, {3, "Whiskers"}}

	q := FromSlice(people).
		SelectMany(func(p person) Query[KeyValue[string, int]] {
			return FromSlice(p.petIDs).Select(func(id int) KeyValue[string, int] {
				return KeyValue[string, int]{p.name, id}
			})
		}).
		Join(
			FromSlice(pets),
			func(kv KeyValue[string, int]) int { return kv.Value },
			func(p pet) int { return p.id },
			func(kv KeyValue[string, int], p pet) string {
				return kv.Key + ":" + p.name
			})

	want := []string{"Ahmet:Fluffy", "Ahmet:Rex", "Bora:Whiskers"}
	if !testQueryIteration(q, want) {
		t.Errorf("Join()=%v expected %v", toSlice(q), want)
	}
}
