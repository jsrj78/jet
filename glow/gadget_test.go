package main

import (
	"bytes"
	"testing"
)

func TestPrintGadgetExists(t *testing.T) {
	f, ok := Registry["print"]
	if !ok {
		t.Errorf("could not find print gadget")
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
	g.In(0, NewMsg("hello"))

	if b.String() != "hello" {
		t.Errorf("debug output is incorrect: %v", b.String())
	}
}
