package glow

import (
	"bytes"
	"testing"
)

func TestUnknownGadget(t *testing.T) {
	g := NewGadget("blah")
	if g != nil {
		t.Errorf("expected nil, got: %T", g)
	}
}

func TestPrintGadgetExists(t *testing.T) {
	g := NewGadget("print")
	if g == nil {
		t.Fatalf("could not find [print] gadget")
	}
	_, ok := g.(*Gadget)
	if !ok {
		t.Errorf("not a gadget: %v", g)
	}
}

func TestPrintGadget(t *testing.T) {
	tmp := Debug
	defer func() { Debug = tmp }()
	b := &bytes.Buffer{}
	Debug = b

	g := NewGadget("print")
	g.Feed(0, NewMsg("hello"))

	if b.String() != "hello" {
		t.Errorf("expected 'hello', got: %v", b)
	}
}

func TestPassGadgetExists(t *testing.T) {
	f, ok := Registry["pass"]
	if !ok {
		t.Fatalf("could not find [pass] gadget")
	}
	g := f(nil)
	_, ok = g.(*Gadget)
	if !ok {
		t.Errorf("not a gadget: %v", g)
	}
}

func TestPassAndPrintGadget(t *testing.T) {
	tmp := Debug
	defer func() { Debug = tmp }()
	var b bytes.Buffer
	Debug = &b

	g1 := NewGadget("pass")
	g2 := NewGadget("print")
	g1.Connect(0, g2, 0) // g1.out[0] => g2.in[0]

	g1.Feed(0, NewMsg("howdy"))

	if b.String() != "howdy" {
		t.Errorf("expected 'howdy', got: %v", b)
	}
}

/*
func TestEmptyCircuit(t *testing.T) {
	g := NewGadget("circuit")
	_, ok := g.(*Circuit)
	if !ok {
		t.Errorf("expected circuit, got %T", g)
	}
}
*/

func TestBuildCircuit(t *testing.T) {
	tmp := Debug
	defer func() { Debug = tmp }()
	b := &bytes.Buffer{}
	Debug = b

	c := new(Circuit)
	g := NewGadget("pass")
	c.Add(g)
	c.Add(NewGadget("print"))
	c.AddWire(0, 0, 1, 0)

	g.Feed(0, NewMsg("bingo"))

	if b.String() != "bingo" {
		t.Errorf("expected 'bingo', got: %v", b)
	}
}

func TestCircuitInlet(t *testing.T) {
	tmp := Debug
	defer func() { Debug = tmp }()
	b := &bytes.Buffer{}
	Debug = b

	c := new(Circuit)
	c.Add(NewGadget("inlet~"))
	c.Add(NewGadget("print"))
	c.AddWire(0, 0, 1, 0)

	c.Feed(0, NewMsg("foo"))

	if b.String() != "foo" {
		t.Errorf("expected 'foo', got: %v", b)
	}
}

func TestCircuitOutlet(t *testing.T) {
	tmp := Debug
	defer func() { Debug = tmp }()
	b := &bytes.Buffer{}
	Debug = b

	c := new(Circuit)
	c.Add(NewGadget("inlet~"))
	c.Add(NewGadget("outlet~"))
	c.AddWire(0, 0, 1, 0)

	g := NewGadget("print")
	c.Connect(0, g, 0) // c.out[0] => g.in[0]

	c.Feed(0, NewMsg("bar"))

	if b.String() != "bar" {
		t.Errorf("expected 'bar', got: %v", b)
	}
}
