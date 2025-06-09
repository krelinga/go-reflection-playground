package testtypes

import (
	"fmt"
	"reflect"
)

type IFace interface {
	String() string
}

type IFaceImpl int

func (i IFaceImpl) String() string {
	return fmt.Sprint(int(i))
}

func NewIFaceValue(i int) reflect.Value {
	return reflect.ValueOf(IFaceImpl(i)).Convert(reflect.TypeFor[IFace]())
}

type Inner struct {
	Int int
}

type Outer struct {
	Inner
}
