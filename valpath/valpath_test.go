package valpath_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/krelinga/go-reflection-playground/valpath"
)

func TestValPath(t *testing.T) {
	testCases := []struct {
		name    string
		path    valpath.Path
		in      reflect.Value
		wantAny any
		wantErr error
	}{
		{
			name:    "empty path on direct value",
			in:      reflect.ValueOf(int(42)),
			wantAny: int(42),
		},
		{
			name:    "deref on direct value",
			path:    valpath.Path{valpath.Deref{}},
			in:      reflect.ValueOf(int(42)),
			wantErr: valpath.ErrTodo,
		},
		{
			name:    "interface on direct value",
			path:    valpath.Path{valpath.Inter{}},
			in:      reflect.ValueOf(int(42)),
			wantErr: valpath.ErrTodo,
		},
		{
			name:    "index on direct value",
			path:    valpath.Path{valpath.Index(0)},
			in:      reflect.ValueOf(int(42)),
			wantErr: valpath.ErrTodo,
		},
		{
			name:    "map key on direct value",
			path:    valpath.Path{valpath.MapKey(reflect.ValueOf(string("key")))},
			in:      reflect.ValueOf(int(42)),
			wantErr: valpath.ErrTodo,
		},
		{
			name:    "map value on direct value",
			path:    valpath.Path{valpath.MapValueOfKey(reflect.ValueOf(string("key")))},
			in:      reflect.ValueOf(int(42)),
			wantErr: valpath.ErrTodo,
		},
		{
			name:    "exported field on direct value",
			path:    valpath.Path{valpath.ExportedField("Int")},
			in:      reflect.ValueOf(int(42)),
			wantErr: valpath.ErrTodo,
		},
		// TODO: start here and add a lot more tests.
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.wantErr == nil && tt.wantAny == nil) || (tt.wantErr != nil && tt.wantAny != nil) {
				t.Fatalf("Test case %s: wantErr and wantAny must not be both nil or both non-nil", tt.name)
			}

			got, err := tt.path.Traverse(tt.in)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Test case %s: got error %v, want %v", tt.name, err, tt.wantErr)
				}
				if got.IsValid() {
					t.Errorf("Test case %s: got value %v, want invalid value", tt.name, got)
				}
			} else {
				if err != nil {
					t.Errorf("Test case %s: got error %v, want no error", tt.name, err)
				}
				if !got.IsValid() {
					t.Errorf("Test case %s: got invalid value, want valid value", tt.name)
				} else if !reflect.DeepEqual(got.Interface(), tt.wantAny) {
					t.Errorf("Test case %s: got value %v, want %v", tt.name, got.Interface(), tt.wantAny)
				}
			}
		})
	}
	// Create a pointer to an int
	i := 42
	ptr := &i

	// Create a reflect.Value from the pointer
	val := reflect.ValueOf(ptr)

	// Create a Deref element
	derefElem := valpath.Deref{}

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
