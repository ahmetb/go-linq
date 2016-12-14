package linq

import "testing"

import "fmt"

type foo struct {
	f1 int
	f2 bool
	f3 string
}

func (f foo) Iterate() Iterator {
	i := 0

	return func() (item interface{}, ok bool) {
		switch i {
		case 0:
			item = f.f1
			ok = true
		case 1:
			item = f.f2
			ok = true
		case 2:
			item = f.f3
			ok = true
		default:
			ok = false
		}

		i++
		return
	}
}

func (f foo) CompareTo(c Comparable) int {
	a, b := f.f1, c.(foo).f1

	if a < b {
		return -1
	} else if a > b {
		return 1
	}

	return 0
}

func toSlice(q Query) (result []interface{}) {
	next := q.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		result = append(result, item)
	}

	return
}

func validateQuery(q Query, output []interface{}) bool {
	next := q.Iterate()

	for _, oitem := range output {
		qitem, _ := next()

		if oitem != qitem {
			return false
		}
	}

	_, ok := next()
	_, ok2 := next()
	return !(ok || ok2)
}

func mustPanicWithError(t *testing.T, expectedErr string, f func()) {
	defer func() {
		r := recover()
		err := fmt.Sprintf("%s", r)
		if err != expectedErr {
			t.Fatalf("got=[%v] expected=[%v]", err, expectedErr)
		}
	}()
	f()
}
