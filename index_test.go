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
