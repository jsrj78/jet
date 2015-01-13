package flow

import "testing"

func TestNewCircuit(t *testing.T) {
	c := NewCircuit()
	if c == nil {
		t.Fail()
	}
	c.Terminate()
}

type test1G struct {
	Gadget
}

func init() {
	Register("test1", func() Circuitry { return new(test1G) })
}

func TestAddTestCircuit(t *testing.T) {
	c := NewCircuit()
	c.Request("add", "test1")
	c.Terminate()
}
