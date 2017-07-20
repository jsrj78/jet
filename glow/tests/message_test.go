package tests

import (
	"testing"

	"github.com/jeelabs/jet/glow"
)

func TestDefaultMessageIsBang(t *testing.T) {
	var m glow.Message
	if !m.IsBang() {
		t.Errorf("should be bang")
	}
}

func TestAtOutOfRangeIsBang(t *testing.T) {
	var m glow.Message
	if !m.At(123).IsBang() {
		t.Errorf("should be bang")
	}
}

func TestIntMessage(t *testing.T) {
	m := glow.Message{123}
	if m.IsBang() {
		t.Errorf("should not be bang")
	}
	if !m.IsInt() {
		t.Errorf("should be int")
	}
	if m[0] != 123 {
		t.Errorf("should be 123")
	}
}

func TestNotIntMessage(t *testing.T) {
	m := glow.Message{123, 456}
	if m.IsInt() {
		t.Errorf("should not be int")
	}
}

func TestStringMessage(t *testing.T) {
	m := glow.Message{"abc"}
	if m.IsBang() {
		t.Errorf("should not be bang")
	}
	if m.IsInt() {
		t.Errorf("should not be int")
	}
	if !m.IsString() {
		t.Errorf("should be string")
	}
	if m[0] != "abc" {
		t.Errorf("should be \"abc\"")
	}
}

func TestNotStringMessage(t *testing.T) {
	m := glow.Message{"abc", "def"}
	if m.IsString() {
		t.Errorf("should not be string")
	}
}

func TestAsInt(t *testing.T) {
	m := glow.Message{12345}
	if m.AsInt() != 12345 {
		t.Errorf("expected 12345, got: %d", m.AsInt())
	}
}

func TestAsNotInt(t *testing.T) {
	m := glow.Message{"abc"}
	if m.AsInt() != 0 {
		t.Errorf("expected 0, got: %d", m.AsInt())
	}
}

func TestAsString(t *testing.T) {
	m := glow.Message{"abcde"}
	if m.AsString() != "abcde" {
		t.Errorf("expected \"abcde\", got: %s", m.AsString())
	}
}

func TestAsNotString(t *testing.T) {
	m := glow.Message{123}
	if m.AsString() != "" {
		t.Errorf("expected \"\", got: %s", m.AsString())
	}
}

var nestedMessage = glow.Message{123, "abc", glow.Message{4, nil, 6}, "d e", 789, "f\ng"}

func TestNestedMessage(t *testing.T) {
	m := nestedMessage
	if len(m) != 6 {
		t.Errorf("expected length 6, got: %d", len(m))
	}
	if len(m.At()) != 6 {
		t.Errorf("should be length 6")
	}
	if x := m.At(2); len(x) != 3 {
		t.Errorf("expected length 3, got: %d", len(x))
	}
	if x := m.At(2, 0); len(x) != 1 {
		t.Errorf("expected length 1, got: %d", len(x))
	}
	if x := m.At(2, 0); !x.IsInt() {
		t.Errorf("expected int, got: %s", x)
	}
	if x := m.At(2, 1); !x.IsBang() {
		t.Errorf("expected bang, got: %s", x)
	}
	if x := m.At(2, 2).AsInt(); x != 6 {
		t.Errorf("expected 6, got: %d", x)
	}
}

func TestMessageAsString(t *testing.T) {
	var s string
	s = glow.Message{}.String()
	if s != "[]" {
		t.Errorf("wrong string, got: %s", s)
	}
	s = glow.Message{"abc"}.String()
	if s != "abc" {
		t.Errorf("wrong string, got: %s", s)
	}
	s = glow.Message{""}.String()
	if s != `""` {
		t.Errorf("wrong string, got: %s", s)
	}
	s = glow.Message{"123"}.String()
	if s != `"123"` {
		t.Errorf("wrong string, got: %s", s)
	}
}

func TestNestedMessageAsString(t *testing.T) {
	s := nestedMessage.String()
	if s != `123 abc [4 [] 6] "d e" 789 "f\ng"` {
		t.Errorf("wrong string, got: %s", s)
	}
}

func TestNilInMessage(t *testing.T) {
	m := glow.Message{nil}
	if m.String() != "[]" {
		t.Error("expected [], got:", m.String())
	}
}
