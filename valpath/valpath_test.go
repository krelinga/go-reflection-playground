package valpath_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/krelinga/go-reflection-playground/testtypes"
	"github.com/krelinga/go-reflection-playground/valpath"
)

func TestValPath(t *testing.T) {
	type Sub struct {
		name    string
		path    valpath.Path
		wantAny any
		wantErr error
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
					name:    "empty path",
					wantAny: int(42),
				},
				{
					name:    "deref",
					path:    valpath.Path{valpath.Deref()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "interface",
					path:    valpath.Path{valpath.Inter()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "index",
					path:    valpath.Path{valpath.Index(0)},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map key",
					path:    valpath.Path{valpath.MapKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map value",
					path:    valpath.Path{valpath.MapValueOfKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "exported field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantErr: valpath.ErrTodo,
				},
			},
		},
		{
			name: "interface value",
			in:   testtypes.NewIFaceValue(42),
			sub: []Sub{
				{
					name:    "empty path",
					wantAny: testtypes.IFaceImpl(42),
				},
				{
					name:    "deref",
					path:    valpath.Path{valpath.Deref()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "interface",
					path:    valpath.Path{valpath.Inter()},
					wantAny: testtypes.IFaceImpl(42),
				},
				{
					name:    "index",
					path:    valpath.Path{valpath.Index(0)},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map key",
					path:    valpath.Path{valpath.MapKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map value",
					path:    valpath.Path{valpath.MapValueOfKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "exported field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantErr: valpath.ErrTodo,
				},
			},
		},
		{
			name: "slice of int",
			in:   reflect.ValueOf([]int{1, 2, 3}),
			sub: []Sub{
				{
					name:    "empty path",
					wantAny: []int{1, 2, 3},
				},
				{
					name:    "deref",
					path:    valpath.Path{valpath.Deref()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "interface",
					path:    valpath.Path{valpath.Inter()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "index",
					path:    valpath.Path{valpath.Index(0)},
					wantAny: int(1),
				},
				{
					name:    "map key",
					path:    valpath.Path{valpath.MapKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map value",
					path:    valpath.Path{valpath.MapValueOfKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "exported field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantErr: valpath.ErrTodo,
				},
			},
		},
		{
			name: "map string to int",
			in:   reflect.ValueOf(map[string]int{"key": 42}),
			sub: []Sub{
				{
					name:    "empty path",
					wantAny: map[string]int{"key": 42},
				},
				{
					name:    "deref",
					path:    valpath.Path{valpath.Deref()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "interface",
					path:    valpath.Path{valpath.Inter()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "index",
					path:    valpath.Path{valpath.Index(0)},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map key",
					path:    valpath.Path{valpath.MapKey("key")},
					wantAny: string("key"),
				},
				{
					name:    "map value",
					path:    valpath.Path{valpath.MapValueOfKey("key")},
					wantAny: int(42),
				},
				{
					name:    "exported field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantErr: valpath.ErrTodo,
				},
			},
		},
		{
			name: "Struct with exported field",
			in:   reflect.ValueOf(struct{ Int int }{Int: 42}),
			sub: []Sub{
				{
					name:    "empty path",
					wantAny: struct{ Int int }{Int: 42},
				},
				{
					name:    "deref",
					path:    valpath.Path{valpath.Deref()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "interface",
					path:    valpath.Path{valpath.Inter()},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "index",
					path:    valpath.Path{valpath.Index(0)},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map key",
					path:    valpath.Path{valpath.MapKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "map value",
					path:    valpath.Path{valpath.MapValueOfKey("key")},
					wantErr: valpath.ErrTodo,
				},
				{
					name:    "exported field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantAny: int(42),
				},
			},
		},
		{
			name: "Struct with exported embedded struct (by value)",
			in:   reflect.ValueOf(testtypes.Outer{Inner: testtypes.Inner{Int: 42}}),
			sub: []Sub{
				// NOTE: this test uses a different set of sub-cases than the others above.
				{
					name:    "access promoted field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantAny: int(42),
				},
				{
					name:    "access non-promoted field",
					path:    valpath.Path{valpath.ExportedField("Inner"), valpath.ExportedField("Int")},
					wantAny: int(42),
				},
			},
		},
		{
			name: "Struct with exported embedded struct (by pointer)",
			in:   reflect.ValueOf(testtypes.OuterPtr{Inner: &testtypes.Inner{Int: 42}}),
			sub: []Sub{
				// NOTE: this test uses a different set of sub-cases than the others above.
				{
					name:    "access promoted field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantAny: int(42),
				},
				{
					name: "access non-promoted field",
					path: valpath.Path{
						valpath.ExportedField("Inner"),
						valpath.Deref(),
						valpath.ExportedField("Int")},
					wantAny: int(42),
				},
			},
		},
		{
			name: "Struct with exported embedded struct (by nil pointer)",
			in:   reflect.ValueOf(testtypes.OuterPtr{Inner: nil}),
			sub: []Sub{
				// NOTE: this test uses a different set of sub-cases than the others above.
				{
					name:    "access promoted field",
					path:    valpath.Path{valpath.ExportedField("Int")},
					wantErr: valpath.ErrTodo,
				},
				{
					name: "access non-promoted field",
					path: valpath.Path{
						valpath.ExportedField("Inner"),
						valpath.Deref(),
						valpath.ExportedField("Int")},
					wantErr: valpath.ErrTodo,
				},
			},
		},
		// TODO: start here and add a lot more tests.
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			for _, sub := range tt.sub {
				t.Run(sub.name, func(t *testing.T) {
					if (sub.wantErr == nil && sub.wantAny == nil) || (sub.wantErr != nil && sub.wantAny != nil) {
						t.Fatal("wantErr and wantAny must not be both nil or both non-nil")
					}

					got, err := sub.path.Traverse(tt.in)
					if sub.wantErr != nil {
						if !errors.Is(err, sub.wantErr) {
							t.Errorf("got error %v, want %v", err, sub.wantErr)
						}
						if got.IsValid() {
							t.Errorf("got value %v, want invalid value", got)
						}
					} else {
						if err != nil {
							t.Errorf("got error %v, want no error", err)
						}
						if !got.IsValid() {
							t.Error("got invalid value, want valid value")
						} else if !reflect.DeepEqual(got.Interface(), sub.wantAny) {
							t.Errorf("got value %v, want %v", got.Interface(), sub.wantAny)
						}
					}
				})
			}
		})
	}
	// Create a pointer to an int
	i := 42
	ptr := &i

	// Create a reflect.Value from the pointer
	val := reflect.ValueOf(ptr)

	// Create a Deref element
	derefElem := valpath.Deref()

	// Traverse the value using Deref
	result, err := derefElem.Traverse(val)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the result is the dereferenced value
	if result.Kind() != reflect.Int || result.Int() != 42 {
		t.Fatalf("expected dereferenced value to be 42, got %v", result)
	}
}
