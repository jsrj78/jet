package glow

import "testing"

func TestDummy(t *testing.T) {
	if true == false {
		t.Errorf("impossible!")
	}
}
