package valpath

import (
	"errors"
	"fmt"
	"iter"
	"reflect"
	"strings"

	"github.com/krelinga/go-iters"
)

var ErrTodo = errors.New("TODO: define a better error for this")

var zeroValue = reflect.Value{}

type Elem interface {
	String() string

	Traverse(reflect.Value) (reflect.Value, error)

	Elems() iter.Seq[Elem]
}

type Path []Elem

func (p Path) String() string {
	b := &strings.Builder{}
	for elem := range p.Elems() {
		if b.Len() > 0 {
			b.WriteString(" / ")
		}
		if elem == nil {
			b.WriteString("<nil>")
			continue
		}
		b.WriteString(elem.String())
	}
	if b.Len() == 0 {
		return "<empty path>"
	} else {
		return b.String()
	}
}

func (p Path) Traverse(v reflect.Value) (reflect.Value, error) {
	for elem := range p.Elems() {
		if val, err := elem.Traverse(v); err != nil {
			return zeroValue, ErrTodo
		} else {
			v = val
		}
	}
	return v, nil
}

func (p Path) Elems() iter.Seq[Elem] {
	children := make([]iter.Seq[Elem], len(p))
	for i, elem := range p {
		children[i] = elem.Elems()
	}
	return iters.Concat(children...)
}

func Deref() Elem {
	return DerefElem{}
}

type DerefElem struct{}

func (d DerefElem) String() string {
	return "<deref>"
}

func (d DerefElem) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	if v.Kind() != reflect.Pointer {
		return zeroValue, ErrTodo
	}
	if v.IsNil() {
		return zeroValue, ErrTodo
	}
	return v.Elem(), nil
}

func (d DerefElem) Elems() iter.Seq[Elem] {
	return func(yield func(Elem) bool) {
		yield(d)
	}
}

func Inter() Elem {
	return InterElem{}
}

type InterElem struct{}

func (i InterElem) String() string {
	return "<inter>"
}

func (i InterElem) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	if v.Kind() != reflect.Interface {
		return zeroValue, ErrTodo
	}
	if v.IsNil() {
		return zeroValue, ErrTodo
	}
	return v.Elem(), nil
}

func (i InterElem) Elems() iter.Seq[Elem] {
	return func(yield func(Elem) bool) {
		yield(i)
	}
}

func Index(i int) Elem {
	return IndexElem(i)
}

type IndexElem int

func (i IndexElem) String() string {
	return fmt.Sprintf("<index %d>", i)
}

func (i IndexElem) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return zeroValue, ErrTodo
	}
	if i < 0 || i >= IndexElem(v.Len()) {
		return zeroValue, ErrTodo
	}
	return v.Index(int(i)), nil
}

func (i IndexElem) Elems() iter.Seq[Elem] {
	return func(yield func(Elem) bool) {
		yield(i)
	}
}

func MapKey[K comparable](k K) Elem {
	return MapKeyElem(reflect.ValueOf(k))
}

type MapKeyElem reflect.Value

func (m MapKeyElem) String() string {
	return fmt.Sprintf("<map key %s>", reflect.Value(m).String())
}

func (m MapKeyElem) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	if v.Kind() != reflect.Map {
		return zeroValue, ErrTodo
	}

	key := reflect.Value(m)
	if !key.IsValid() {
		return zeroValue, ErrTodo
	}
	if !key.Type().AssignableTo(v.Type().Key()) {
		return zeroValue, ErrTodo
	}

	if v.IsNil() {
		return zeroValue, ErrTodo
	}
	found := v.MapIndex(key)
	if !found.IsValid() {
		return zeroValue, ErrTodo
	}

	return key, nil
}

func (m MapKeyElem) Elems() iter.Seq[Elem] {
	return func(yield func(Elem) bool) {
		yield(m)
	}
}

func MapValueOfKey[K comparable](k K) Elem {
	return MapValueOfKeyElem(reflect.ValueOf(k))
}

type MapValueOfKeyElem reflect.Value

func (m MapValueOfKeyElem) String() string {
	return fmt.Sprintf("<map value of key %s>", reflect.Value(m).String())
}

func (m MapValueOfKeyElem) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	if v.Kind() != reflect.Map {
		return zeroValue, ErrTodo
	}

	key := reflect.Value(m)
	if !key.IsValid() {
		return zeroValue, ErrTodo
	}
	if !key.Type().AssignableTo(v.Type().Key()) {
		return zeroValue, ErrTodo
	}

	if v.IsNil() {
		return zeroValue, ErrTodo
	}
	val := v.MapIndex(key)
	if !val.IsValid() {
		return zeroValue, ErrTodo
	}

	return val, nil
}

func (m MapValueOfKeyElem) Elems() iter.Seq[Elem] {
	return func(yield func(Elem) bool) {
		yield(m)
	}
}

func ExportedField(name string) Elem {
	return ExportedFieldElem(name)
}

type ExportedFieldElem string

func (f ExportedFieldElem) String() string {
	return fmt.Sprintf("<exported field %s>", string(f))
}

// TODO: currently this supports finding promoted fields from embedded structs.  It isn't clear to me that
// this is actually something we want to support.  It would mean that there is more than one way to
// address a field, which seems like it could lead to confusion.
func (f ExportedFieldElem) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	if v.Kind() != reflect.Struct {
		return zeroValue, ErrTodo
	}
	t := v.Type()
	fieldDesc, ok := t.FieldByName(string(f))
	if !ok {
		return zeroValue, ErrTodo
	}
	if !fieldDesc.IsExported() {
		return zeroValue, ErrTodo
	}

	fieldValue, err := v.FieldByIndexErr(fieldDesc.Index)
	if err != nil {
		// This happens if the field requires traversing a nil pointer.
		return zeroValue, ErrTodo
	}
	return fieldValue, nil
}

func (f ExportedFieldElem) Elems() iter.Seq[Elem] {
	return func(yield func(Elem) bool) {
		yield(f)
	}
}
