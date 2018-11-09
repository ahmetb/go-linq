package linq

import (
	"testing"
)

func TestLinqLeftJoinWithNullMarkerSample(t *testing.T) {

	type str1 struct {
		col1 string
		col2 string
	}

	type str2 struct {
		col1 string
		col2 string
	}

	type joinedData struct {
		part1     str1
		part2     str2
		part2null bool
	}

	data1 := make([]str1, 0)
	data1 = append(data1, str1{"a", "1"})
	data1 = append(data1, str1{"b", "2"})
	data1 = append(data1, str1{"c", "3"})
	data2 := make([]str2, 0)
	data2 = append(data2, str2{"a", "1"})
	data2 = append(data2, str2{"a", "1"})
	data2 = append(data2, str2{"a", "4"})

	result := From(data1).LeftJoin(
		From(data2),
		func(outer interface{}) interface{} { return outer.(str1).col1 },
		func(inner interface{}) interface{} { return inner.(str2).col1 },
		func(outer interface{}, innner interface{}) interface{} {
			return joinedData{outer.(str1), innner.(str2), false}
		},
		func(outer interface{}) interface{} {
			return joinedData{outer.(str1), str2{}, true}
		})
	want := []interface{}{
		joinedData{str1{"a", "1"}, str2{"a", "1"}, false},
		joinedData{str1{"a", "1"}, str2{"a", "1"}, false},
		joinedData{str1{"a", "1"}, str2{"a", "4"}, false},
		joinedData{str1{"b", "2"}, str2{"", ""}, true},
		joinedData{str1{"c", "3"}, str2{"", ""}, true},
	}

	if !validateQuery(result, want) {
		t.Errorf("From().Join()=%v expected %v", toSlice(result), want)
	}
}

func TestLeftJoinExpectedValues(t *testing.T) {
	outer := []int{0, 1, 2, 3, 4, 5, 8}
	inner := []int{1, 2, 1, 4, 7, 6, 7, 2}
	want := []interface{}{
		KeyValue{0, -1},
		KeyValue{1, 1},
		KeyValue{1, 1},
		KeyValue{2, 2},
		KeyValue{2, 2},
		KeyValue{3, -1},
		KeyValue{4, 4},
		KeyValue{5, -1},
		KeyValue{8, -1},
	}

	q := From(outer).LeftJoin(
		From(inner),
		func(i interface{}) interface{} { return i },
		func(i interface{}) interface{} { return i },
		func(outer interface{}, inner interface{}) interface{} {
			return KeyValue{outer, inner}
		},
		func(outer interface{}) interface{} {
			return KeyValue{outer, -1}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().Join()=%v expected %v", toSlice(q), want)
	}
}
