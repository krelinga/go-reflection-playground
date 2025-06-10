package valpattern_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-iters"
	"github.com/krelinga/go-reflection-playground/valpath"
	"github.com/krelinga/go-reflection-playground/valpattern"
	"github.com/krelinga/go-sets"
)

func checkEqual(t *testing.T, got, want []iters.Pair[valpath.Path, reflect.Value]) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got), len(want))
		return
	}
	claimed := sets.New[int]()
wants:
	for _, w := range want {
		for i, g := range got {
			if !claimed.Has(i) && reflect.DeepEqual(g, w) {
				claimed.Add(i)
				continue wants
			}
		}
	}
	if claimed.Len() == len(want) {
		return
	}
	t.Errorf("mismatch: got %v, want %v", got, want)
}

func TestPattern(t *testing.T) {
	type Sub struct {
		name    string
		pattern valpattern.Pattern
		want    []iters.Pair[valpath.Path, reflect.Value]
	}
	testCases := []struct {
		name string
		in   reflect.Value
		sub  []Sub
	}{
		{
			name: "direct int value",
			in:   reflect.ValueOf(int(42)),
			sub: []Sub{
				{
					name:    "empty pattern",
					pattern: valpattern.Empty(),
					want: []iters.Pair[valpath.Path, reflect.Value]{
						iters.NewPair(valpath.Empty(), reflect.ValueOf(int(42))),
					},
				},
				// TODO: add many more tests for other kinds of patterns, and more input values too.
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			for _, sub := range tt.sub {
				t.Run(sub.name, func(t *testing.T) {
					got := slices.Collect(iters.ToPairs(sub.pattern.Match(tt.in)))
					checkEqual(t, got, sub.want)
				})
			}
		})
	}
}
