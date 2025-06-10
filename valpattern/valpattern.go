package valpattern

import (
	"iter"
	"reflect"
	"slices"
	"strings"

	"github.com/krelinga/go-iters"
	"github.com/krelinga/go-reflection-playground/valpath"
)

type Pattern interface {
	String() string
	Match(reflect.Value) iter.Seq2[valpath.Path, reflect.Value]
	elems() iter.Seq[Pattern]
}

func Path(p valpath.Path) Pattern {
	return pathPat{p}
}

type pathPat struct {
	Path valpath.Path
}

func (p pathPat) String() string {
	return p.Path.String()
}

func (p pathPat) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if found, err := p.Path.Traverse(v); err != nil {
		return iters.Empty2[valpath.Path, reflect.Value]()
	} else {
		return iters.Single2((p.Path), found)
	}
}

func (p pathPat) elems() iter.Seq[Pattern] {
	return iters.Single(Pattern(p))
}

func AllExportedFields() Pattern {
	return allExportedFieldsPat{}
}

type allExportedFieldsPat struct{}

func (allExportedFieldsPat) String() string {
	return "<all exported fields>"
}

func (allExportedFieldsPat) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return iters.Empty2[valpath.Path, reflect.Value]()
	}
	allFields := slices.Values(reflect.VisibleFields(v.Type()))
	exportedFields := iters.Filter(allFields, reflect.StructField.IsExported)
	pairs := iters.Map(exportedFields, func(f reflect.StructField) iters.Pair[valpath.Path, reflect.Value] {
		return iters.NewPair(valpath.ExportedField(f.Name), v.FieldByName(f.Name))
	})
	return iters.FromPairs(pairs)
}

func (allExportedFieldsPat) elems() iter.Seq[Pattern] {
	return iters.Single(AllExportedFields())
}

func AllMapKeys() Pattern {
	return allMapKeysPat{}
}

type allMapKeysPat struct{}

func (allMapKeysPat) String() string {
	return "<all map keys>"
}

func (allMapKeysPat) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if !v.IsValid() || v.Kind() != reflect.Map || v.IsNil() {
		return iters.Empty2[valpath.Path, reflect.Value]()
	}

	keys := slices.Values(v.MapKeys())
	pairs := iters.Map(keys, func(k reflect.Value) iters.Pair[valpath.Path, reflect.Value] {
		return iters.NewPair(valpath.MapKey(k), v.MapIndex(k))
	})
	return iters.FromPairs(pairs)
}

func (allMapKeysPat) elems() iter.Seq[Pattern] {
	return iters.Single(AllMapKeys())
}

func AllMapValues() Pattern {
	return allMapValuesPat{}
}

type allMapValuesPat struct{}

func (allMapValuesPat) String() string {
	return "<all map values>"
}

func (allMapValuesPat) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if !v.IsValid() || v.Kind() != reflect.Map || v.IsNil() {
		return iters.Empty2[valpath.Path, reflect.Value]()
	}

	entries := func(yield func(k, v reflect.Value) bool) {
		mapRange := v.MapRange()
		for mapRange.Next() {
			if !yield(mapRange.Key(), mapRange.Value()) {
				return
			}
		}
	}
	return iters.Map2(entries, func(k, v reflect.Value) (valpath.Path, reflect.Value) {
		return valpath.MapValueOfKey(k), v
	})
}

func (allMapValuesPat) elems() iter.Seq[Pattern] {
	return iters.Single(AllMapValues())
}

func Join(children ...Pattern) Pattern {
	asIter := slices.Values(children)
	nonNil := iters.Filter(asIter, func(e Pattern) bool {
		return e != nil
	})
	childrenIters := iters.Map(nonNil, Pattern.elems)
	children = slices.Collect(iters.Concat(slices.Collect(childrenIters)...))

	switch len(children) {
	case 0:
		return emptyPat{}
	case 1:
		return children[0]
	default:
		return joinedPat(children)
	}
}

type joinedPat []Pattern

func (j joinedPat) String() string {
	b := &strings.Builder{}
	for elem := range j.elems() {
		if b.Len() > 0 {
			b.WriteString(" / ")
		}
		b.WriteString(elem.String())
	}
	return b.String()
}

func (j joinedPat) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if !v.IsValid() {
		return iters.Empty2[valpath.Path, reflect.Value]()
	}

	out := slices.Collect(iters.ToPairs(Empty().Match(v)))
	for elem := range j.elems() {
		existing := slices.Values(out)
		newChildren := iters.Map(existing, func(in iters.Pair[valpath.Path, reflect.Value]) iter.Seq[iters.Pair[valpath.Path, reflect.Value]] {
			oldPath := in.One
			oldVal := in.Two
			matches := elem.Match(oldVal)
			withFixedPath := iters.Map2(matches, func(p valpath.Path, v reflect.Value) (valpath.Path, reflect.Value) {
				return valpath.Join(oldPath, p), v
			})
			return iters.ToPairs(withFixedPath)
		})
		flattened := iters.Concat(slices.Collect(newChildren)...)
		out = slices.Collect(flattened)
	}
	return iters.FromPairs(slices.Values(out))
}

func (j joinedPat) elems() iter.Seq[Pattern] {
	children := make([]iter.Seq[Pattern], len(j))
	for i, elem := range j {
		children[i] = elem.elems()
	}
	return iters.Concat(children...)
}

type emptyPat struct{}

func (emptyPat) String() string {
	return "<empty pattern>"
}

func (emptyPat) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if !v.IsValid() {
		return iters.Empty2[valpath.Path, reflect.Value]()
	}
	return iters.Single2(valpath.Empty(), v)
}

func (emptyPat) elems() iter.Seq[Pattern] {
	return iters.Empty[Pattern]()
}

func Empty() Pattern {
	return emptyPat{}
}
