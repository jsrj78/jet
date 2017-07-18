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
		t.Errorf("expected \"hello\", got: %v", b)
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
		t.Errorf("expected \"howdy\", got: %v", b)
	}
}
