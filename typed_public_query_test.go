package linq

import (
	"slices"
	"testing"
)

func TestPublicQueryAliasExposesFluentTypedMethods(t *testing.T) {
	var source Query[int] = fromSlice([]int{1, 2, 3})
	got := source.Where(func(value int) bool { return value > 1 }).Select(func(value int) string {
		return string(rune('a' + value))
	}).Results()

	if !slices.Equal(got, []string{"c", "d"}) {
		t.Fatalf("public Query pipeline = %v, want [c d]", got)
	}
}
