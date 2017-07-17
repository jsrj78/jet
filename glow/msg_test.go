package main

import "testing"

func TestDefaultMsgIsBang(t *testing.T) {
	m := NewMsg()
	if !m.IsBang() {
		t.Errorf("should be bang")
	}
}

func TestAtOutOfRangeIsBang(t *testing.T) {
	m := NewMsg()
	if !m.At(123).IsBang() {
		t.Errorf("should be bang")
	}
}

func TestIntMsg(t *testing.T) {
	m := NewMsg(123)
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

func TestNotIntMsg(t *testing.T) {
	m := NewMsg(123, 456)
	if m.IsInt() {
		t.Errorf("should not be int")
	}
}

func TestStringMsg(t *testing.T) {
	m := NewMsg("abc")
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

func TestNotStringMsg(t *testing.T) {
	m := NewMsg("abc", "def")
	if m.IsString() {
		t.Errorf("should not be string")
	}
}

var nestedMsg = NewMsg(123, "abc", NewMsg(4, NewMsg(), 6), "def", 789)

func TestNestedMsg(t *testing.T) {
	m := nestedMsg
	if len(m) != 5 {
		t.Errorf("expected length 4, got %d", len(m))
	}
	m2 := m.At(2)
	if len(m2) != 3 {
		t.Errorf("expected length 3, got %d", len(m2))
	}
}

func TestAsInt(t *testing.T) {
	m := NewMsg(12345)
	if m.AsInt() != 12345 {
		t.Errorf("expected 12345, got %d", m.AsInt())
	}
}

func TestAsNotInt(t *testing.T) {
	m := NewMsg("abc")
	if m.AsInt() != 0 {
		t.Errorf("expected 0, got %d", m.AsInt())
	}
}

func TestAsString(t *testing.T) {
	m := NewMsg("abcde")
	if m.AsString() != "abcde" {
		t.Errorf("expected \"abcde\", got %s", m.AsString())
	}
}

func TestAsNotString(t *testing.T) {
	m := NewMsg(123)
	if m.AsString() != "" {
		t.Errorf("expected \"\", got %s", m.AsString())
	}
}

// At could be dropped, add optional args to AsInt, etc
// AsMsg can return self of any nested element, arbitrarily deep
// IsInt etc could also support nested indexing
