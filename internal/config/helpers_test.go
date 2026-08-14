package config_test

import (
	"reflect"
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

// assertNoNilFields recursively verifies that every pointer, slice, map, and
// interface field reachable from v is non-nil. Scalar and skipped (unexported)
// fields are ignored. On failure it reports the full path, e.g. Keys.Main.Quit.
func assertNoNilFields(t *testing.T, v any) {
	t.Helper()
	var check func(string, reflect.Value)
	check = func(path string, rv reflect.Value) {
		switch rv.Kind() {
		case reflect.Pointer, reflect.Interface:
			if rv.IsNil() {
				t.Errorf("%s is nil", path)
				return
			}
			check(path, rv.Elem())
		case reflect.Struct:
			for i := range rv.NumField() {
				f := rv.Type().Field(i)
				if !f.IsExported() {
					continue
				}
				check(path+"."+f.Name, rv.Field(i))
			}
		case reflect.Slice, reflect.Map:
			if rv.IsNil() {
				t.Errorf("%s is nil", path)
			}
		default:
			return
		}
	}
	check("", reflect.ValueOf(v))
}

func assertKeyBinding(t *testing.T, name string, b *config.KeyBinding, want ...string) {
	t.Helper()
	if b == nil {
		t.Errorf("%s: nil binding, want %v", name, want)
		return
	}
	if got := []string(*b); !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func assertNilKeyBinding(t *testing.T, name string, b *config.KeyBinding) {
	t.Helper()
	if b != nil {
		t.Errorf("%s = %v, want nil", name, *b)
	}
}
