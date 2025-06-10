package valpattern

import (
	"iter"
	"reflect"
	"slices"
	"strings"

	"github.com/krelinga/go-iters"
	"github.com/krelinga/go-reflection-playground/valpath"
)

type Elem interface {
	String() string
	Match(reflect.Value) iter.Seq2[valpath.Path, reflect.Value]
	Elems() iter.Seq[Elem]
}

func Path(p valpath.Path) PathElem {
	return PathElem{p}
}

type PathElem struct {
	Path valpath.Path
}

func (p PathElem) String() string {
	return p.Path.String()
}

func (p PathElem) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if found, err := p.Path.Traverse(v); err != nil {
		return iters.Empty2[valpath.Path, reflect.Value]()
	} else {
		return iters.Single2((p.Path), found)
	}
}

func (p PathElem) Elems() iter.Seq[Elem] {
	return iters.Single(Elem(p))
}

func AllExportedFields() Elem {
	return AllExportedFieldsElem{}
}

type AllExportedFieldsElem struct{}

func (AllExportedFieldsElem) String() string {
	return "<all exported fields>"
}

func (AllExportedFieldsElem) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
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

func (AllExportedFieldsElem) Elems() iter.Seq[Elem] {
	return iters.Single(AllExportedFields())
}

func AllMapKeys() Elem {
	return AllMapKeysElem{}
}

type AllMapKeysElem struct{}

func (AllMapKeysElem) String() string {
	return "<all map keys>"
}

func (AllMapKeysElem) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if !v.IsValid() || v.Kind() != reflect.Map || v.IsNil() {
		return iters.Empty2[valpath.Path, reflect.Value]()
	}

	keys := slices.Values(v.MapKeys())
	pairs := iters.Map(keys, func(k reflect.Value) iters.Pair[valpath.Path, reflect.Value] {
		return iters.NewPair(valpath.MapKey(k), v.MapIndex(k))
	})
	return iters.FromPairs(pairs)
}

func (AllMapKeysElem) Elems() iter.Seq[Elem] {
	return iters.Single(AllMapKeys())
}

func AllMapValues() Elem {
	return AllMapValuesElem{}
}

type AllMapValuesElem struct{}

func (AllMapValuesElem) String() string {
	return "<all map values>"
}

func (AllMapValuesElem) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
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

func (AllMapValuesElem) Elems() iter.Seq[Elem] {
	return iters.Single(AllMapValues())
}

func New(elems ...Elem) Elem {
	asIter := slices.Values(elems)
	nonNil := iters.Filter(asIter, func(e Elem) bool {
		return e != nil
	})
	elems = slices.Collect(nonNil)

	if len(elems) == 0 {
		return nil
	}
	if len(elems) == 1 {
		return elems[0]
	}
	return joined(elems)
}

type joined []Elem

func (j joined) String() string {
	b := &strings.Builder{}
	for elem := range j.Elems() {
		if b.Len() > 0 {
			b.WriteString(" / ")
		}
		b.WriteString(elem.String())
	}
	if b.Len() == 0 {
		return "<empty pattern>"
	} else {
		return b.String()
	}
}

// TODO: think real hard about whether or not this is correct... The logic is very fancy (but very concise?).
func (j joined) Match(v reflect.Value) iter.Seq2[valpath.Path, reflect.Value] {
	if len(j) == 0 {
		// TODO: nil isn't a valid value here, decide how to handle this.
		return iters.Single2(valpath.Path(nil), v)
	}
	if !v.IsValid() {
		return iters.Empty2[valpath.Path, reflect.Value]()
	}

	out := []iters.Pair[valpath.Path, reflect.Value]{}
	for elem := range j.Elems() {
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

func (j joined) Elems() iter.Seq[Elem] {
	children := make([]iter.Seq[Elem], len(j))
	for i, elem := range j {
		children[i] = elem.Elems()
	}
	return iters.Concat(children...)
}
