package main

import (
	"bytes"
	"testing"
)

func TestPrintGadgetExists(t *testing.T) {
	f, ok := Registry["print"]
	if !ok {
		t.Errorf("could not find [print] gadget")
	}
	g := f()
	_, ok = g.(*Gadget)
	if !ok {
		t.Errorf("not a gadget: %v", g)
	}
}

func TestPrintGadget(t *testing.T) {
	tmp := Debug
	defer func() { Debug = tmp }()
	b := &bytes.Buffer{}
	Debug = b

	g := Registry["print"]()
	g.Feed(0, NewMsg("hello"))

	if b.String() != "hello" {
		t.Errorf("expected 'hello', got: %v", b)
	}
}

func TestPassGadgetExists(t *testing.T) {
	f, ok := Registry["pass"]
	if !ok {
		t.Errorf("could not find [pass] gadget")
	}
	g := f()
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

	g1 := Registry["pass"]()
	g2 := Registry["print"]()
	g1.Connect(0, g2, 0) // g1.out[0] => g2.in[0]

	g1.Feed(0, NewMsg("howdy"))

	if b.String() != "howdy" {
		t.Errorf("expected 'howdy', got: %v", b)
	}
}

func TestEmptyCircuit(t *testing.T) {
	g := Registry["circuit"]()
	_, ok := g.(*Circuit)
	if !ok {
		t.Errorf("expected circuit, got %T", g)
	}
}

func TestBuildCircuit(t *testing.T) {
	tmp := Debug
	defer func() { Debug = tmp }()
	b := &bytes.Buffer{}
	Debug = b

	c := new(Circuit)
	g := Registry["pass"]()
	c.Add(g)
	c.Add(Registry["print"]())
	c.Wire(0, 0, 1, 0)

	g.Feed(0, NewMsg("bingo"))

	if b.String() != "bingo" {
		t.Errorf("expected 'bingo', got: %v", b)
	}
}
