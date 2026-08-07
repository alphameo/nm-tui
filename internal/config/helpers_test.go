package config

import (
	"reflect"
	"testing"
)

func assertKeyBinding(t *testing.T, name string, b *KeyBinding, want ...string) {
	t.Helper()
	if b == nil {
		t.Errorf("%s: nil binding, want %v", name, want)
		return
	}
	if got := []string(*b); !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func assertNilKeyBinding(t *testing.T, name string, b *KeyBinding) {
	t.Helper()
	if b != nil {
		t.Errorf("%s = %v, want nil", name, *b)
	}
}
