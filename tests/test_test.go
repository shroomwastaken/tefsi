package tests

import "testing"

func TestTesting(t *testing.T) {
	if 1 != 1 {
		t.Error("one is not one")
	}
}
