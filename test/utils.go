package test

import (
	"testing"
	"runtime/debug"
)

func ExpectSuccess(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected success, got error (%v). %v", err, string(debug.Stack()))
	}
}

func ExpectString(expected, got string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected '%s', got '%s'. %v", expected, got, string(debug.Stack()))
	}
}

func ExpectInt(expected, got int, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %d, got %d. %v", expected, got, string(debug.Stack()))
	}
}

func ExpectBool(expected, got bool, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %v, got %v. %v", expected, got, string(debug.Stack()))
	}
}
