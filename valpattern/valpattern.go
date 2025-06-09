package valpattern

import (
	"iter"
	"reflect"
	"slices"

	"github.com/krelinga/go-iters"
	"github.com/krelinga/go-reflection-playground/valpath"
)

var zeroValue = reflect.Value{}

type Elem interface {
	String() string
	Match(reflect.Value) iter.Seq2[valpath.Elem, reflect.Value]
	Elems() iter.Seq[Elem]
}

func Path(p valpath.Elem) PathElem {
	return PathElem{p}
}

type PathElem struct {
	Path valpath.Elem
}

func (p PathElem) String() string {
	return p.Path.String()
}

func (p PathElem) Match(v reflect.Value) iter.Seq2[valpath.Elem, reflect.Value] {
	if found, err := p.Path.Traverse(v); err != nil {
		return iters.Empty2[valpath.Elem, reflect.Value]()
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

func (AllExportedFieldsElem) Match(v reflect.Value) iter.Seq2[valpath.Elem, reflect.Value] {
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return iters.Empty2[valpath.Elem, reflect.Value]()
	}
	allFields := slices.Values(reflect.VisibleFields(v.Type()))
	exportedFields := iters.Filter(allFields, func(f reflect.StructField) bool {
		return f.IsExported()
	})
	pairs := iters.Map(exportedFields, func(f reflect.StructField) iters.Pair[valpath.Elem, reflect.Value] {
		return iters.NewPair(valpath.ExportedField(f.Name), v.FieldByName(f.Name))
	})
	return iters.FromPairs(pairs)
}

func (AllExportedFieldsElem) Elems() iter.Seq[Elem] {
	return iters.Single(AllExportedFields())
}