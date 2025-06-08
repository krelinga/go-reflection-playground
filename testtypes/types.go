package testtypes

import "fmt"

type T1 int

func (t T1) String() string {
	return fmt.Sprint(int(t))
}

type T2 []int

type T3 map[string]int

type T4 struct {
	T1 T1
	T2 T2
	T3 T3
}

type T5 struct {
	T4
}

type T6 struct {
	Int    int
	String string
}

type T7 struct {
	Float64 float64
	T6      T6
}

type T8 map[T7]T7

type T9 interface {
	String() string
}
