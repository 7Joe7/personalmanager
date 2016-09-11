package test

import (
	"testing"
	"runtime/debug"
)

func ExpectSuccess(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected success, got error (%v). Stack: %v", err, string(debug.Stack()))
	}
}

func ExpectString(expected, got string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected '%s', got '%s'. %v", expected, got, debug.Stack())
	}
}
