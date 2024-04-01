package assert

import (
	"fmt"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got %v; expected: nil", actual)
	}
}

func SliceEqual[T comparable](t *testing.T, actual, expected []T) {
	t.Helper()
	if len(actual) != len(expected) {
		fmt.Println(actual)
		fmt.Println(expected)
		t.Errorf("got: %v; want %v", actual, expected)
	}

	for i, v := range actual {
		if v != expected[i] {
			t.Errorf("got: %v; want %v", actual, expected)
		}
	}
}
