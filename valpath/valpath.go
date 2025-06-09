package valpath

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrTodo = errors.New("TODO: define a better error for this")

var zeroValue = reflect.Value{}

type Element interface {
	String() string

	Traverse(reflect.Value) (reflect.Value, error)

	isElement() // Marker method to identify Element types
}

type Path []Element

func (p Path) String() string {
	if len(p) == 0 {
		return "<empty path>"
	}
	b := &strings.Builder{}
	for i, elem := range p {
		if i > 0 {
			b.WriteString(" / ")
		}
		b.WriteString(elem.String())
	}
	return b.String()
}

func (p Path) Traverse(v reflect.Value) (reflect.Value, error) {
	for _, elem := range p {
		if val, err := elem.Traverse(v); err != nil {
			return zeroValue, ErrTodo
		} else {
			v = val
		}
	}
	return v, nil
}

type Deref struct{}

func (d Deref) String() string {
	return "<deref>"
}

func (d Deref) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (d Deref) isElement() {}

type Inter struct{}

func (i Inter) String() string {
	return "<inter>"
}

func (i Inter) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (i Inter) isElement() {}

type Index int

func (i Index) String() string {
	return fmt.Sprintf("<index %d>", i)
}

func (i Index) Traverse(v reflect.Value) (reflect.Value, error) {
	if !v.IsValid() {
		return zeroValue, ErrTodo
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return zeroValue, ErrTodo
	}
	if i < 0 || i >= Index(v.Len()) {
		return zeroValue, ErrTodo
	}
	return v.Index(int(i)), nil
}

func (i Index) isElement() {}

type MapKey reflect.Value

func (m MapKey) String() string {
	return fmt.Sprintf("<map key %s>", reflect.Value(m).String())
}

func (m MapKey) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (m MapKey) isElement() {}

type MapValueOfKey reflect.Value

func (m MapValueOfKey) String() string {
	return fmt.Sprintf("<map value of key %s>", reflect.Value(m).String())
}

func (m MapValueOfKey) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (m MapValueOfKey) isElement() {}

type ExportedField string

func (f ExportedField) String() string {
	return fmt.Sprintf("<exported field %s>", string(f))
}

// TODO: currently this supports finding promoted fields from embedded structs.  It isn't clear to me that
// this is actually something we want to support.  It would mean that there is more than one way to
// address a field, which seems like it could lead to confusion.
func (f ExportedField) Traverse(v reflect.Value) (reflect.Value, error) {
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

func (f ExportedField) isElement() {}
