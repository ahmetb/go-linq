package linq_test

import (
	"slices"
	"testing"

	linq "github.com/ahmetb/go-linq/v5"
)

type externalWidget struct {
	value int
}

func (widget externalWidget) Value() int { return widget.value }

func TestPublicAPIInfersGenericMethodsAcrossPackageBoundary(t *testing.T) {
	got := linq.FromSlice([]externalWidget{{1}, {2}, {3}}).
		Where(func(widget externalWidget) bool { return widget.value > 1 }).
		Select(externalWidget.Value).
		Results()

	if !slices.Equal(got, []int{2, 3}) {
		t.Fatalf("external generic pipeline = %v, want [2 3]", got)
	}
}

func TestPublicAPIExposesTypedIteration(t *testing.T) {
	var got []int
	for value := range linq.Range(2, 3).Iterate {
		got = append(got, value)
	}

	if !slices.Equal(got, []int{2, 3, 4}) {
		t.Fatalf("external Iterate = %v, want [2 3 4]", got)
	}
}

func TestPublicGroupingAndOrderingTypesCrossPackageBoundary(t *testing.T) {
	var groups linq.Query[linq.Group[int, int]] = linq.FromSlice([]int{3, 1, 2, 4}).GroupBy(
		func(value int) int { return value % 2 },
		func(value int) int { return value },
	)
	grouped := groups.Results()
	if len(grouped) != 2 || !slices.Equal(grouped[0].Group, []int{3, 1}) {
		t.Fatalf("external GroupBy = %v", grouped)
	}

	var ordered linq.OrderedQuery[externalWidget] = linq.FromSlice([]externalWidget{{2}, {1}}).
		OrderBy(externalWidget.Value)
	if got := ordered.Results(); got[0].value != 1 {
		t.Fatalf("external OrderBy = %v", got)
	}
}
