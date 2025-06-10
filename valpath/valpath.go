package valpath

import (
	"errors"
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"

	"github.com/krelinga/go-iters"
)

var ErrTodo = errors.New("TODO: define a better error for this")

var zeroValue = reflect.Value{}

type Path interface {
	String() string

	Traverse(reflect.Value) (reflect.Value, error)

	elems() iter.Seq[Path]
}

func Join(children ...Path) Path {
	asIter := slices.Values(children)
	notNil := iters.Filter(asIter, func(e Path) bool {
		return e != nil
	})
	childLists := iters.Map(notNil, Path.elems)
	children = slices.Collect(iters.Concat(slices.Collect(childLists)...))

	switch len(children) {
	case 0:
		return Empty()
	case 1:
		return children[0]
	default:
		return pathListElem(children)
	}
}

func Empty() Path {
	return emptyPathElem{}
}

type emptyPathElem struct{}

func (e emptyPathElem) String() string {
	return "<empty path>"
}

func (e emptyPathElem) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	return v, nil
}

func (e emptyPathElem) elems() iter.Seq[Path] {
	return iters.Empty[Path]()
}

type pathListElem []Path

func (p pathListElem) String() string {
	b := &strings.Builder{}
	for elem := range p.elems() {
		if b.Len() > 0 {
			b.WriteString(" / ")
		}
		if elem == nil {
			b.WriteString("<nil>")
			continue
		}
		b.WriteString(elem.String())
	}
	return b.String()
}

func (p pathListElem) Traverse(v reflect.Value) (reflect.Value, error) {
	for elem := range p.elems() {
		if val, err := elem.Traverse(v); err != nil {
			return zeroValue, ErrTodo
		} else {
			v = val
		}
	}
	return v, nil
}

func (p pathListElem) elems() iter.Seq[Path] {
	children := make([]iter.Seq[Path], len(p))
	for i, elem := range p {
		children[i] = elem.elems()
	}
	return iters.Concat(children...)
}

func Deref() Path {
	return DerefPart{}
}

type DerefPart struct{}

func (d DerefPart) String() string {
	return "<deref>"
}

func (d DerefPart) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (d DerefPart) elems() iter.Seq[Path] {
	return func(yield func(Path) bool) {
		yield(d)
	}
}

func Inter() Path {
	return InterPart{}
}

type InterPart struct{}

func (i InterPart) String() string {
	return "<inter>"
}

func (i InterPart) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (i InterPart) elems() iter.Seq[Path] {
	return func(yield func(Path) bool) {
		yield(i)
	}
}

func Index(i int) Path {
	return IndexPart(i)
}

type IndexPart int

func (i IndexPart) String() string {
	return fmt.Sprintf("<index %d>", i)
}

func (i IndexPart) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return zeroValue, ErrTodo
	}
	if i < 0 || i >= IndexPart(v.Len()) {
		return zeroValue, ErrTodo
	}
	return v.Index(int(i)), nil
}

func (i IndexPart) elems() iter.Seq[Path] {
	return func(yield func(Path) bool) {
		yield(i)
	}
}

func MapKey[K comparable](k K) Path {
	return MapKeyPart(reflect.ValueOf(k))
}

type MapKeyPart reflect.Value

func (m MapKeyPart) String() string {
	return fmt.Sprintf("<map key %s>", reflect.Value(m).String())
}

func (m MapKeyPart) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (m MapKeyPart) elems() iter.Seq[Path] {
	return func(yield func(Path) bool) {
		yield(m)
	}
}

func MapValueOfKey[K comparable](k K) Path {
	return MapValueOfKeyPart(reflect.ValueOf(k))
}

type MapValueOfKeyPart reflect.Value

func (m MapValueOfKeyPart) String() string {
	return fmt.Sprintf("<map value of key %s>", reflect.Value(m).String())
}

func (m MapValueOfKeyPart) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (m MapValueOfKeyPart) elems() iter.Seq[Path] {
	return func(yield func(Path) bool) {
		yield(m)
	}
}

func ExportedField(name string) Path {
	return ExportedFieldPart(name)
}

type ExportedFieldPart string

func (f ExportedFieldPart) String() string {
	return fmt.Sprintf("<exported field %s>", string(f))
}

// TODO: currently this supports finding promoted fields from embedded structs.  It isn't clear to me that
// this is actually something we want to support.  It would mean that there is more than one way to
// address a field, which seems like it could lead to confusion.
func (f ExportedFieldPart) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (f ExportedFieldPart) elems() iter.Seq[Path] {
	return func(yield func(Path) bool) {
		yield(f)
	}
}
