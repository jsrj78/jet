package gadgets

import (
	"bytes"
	"testing"

	"github.com/jeelabs/jet/glow"
)

func TestUnknownGadget(t *testing.T) {
	g := glow.NewGadget("blah")
	if g != nil {
		t.Errorf("expected nil, got: %T", g)
	}
}

func TestPrintGadgetExists(t *testing.T) {
	g := glow.NewGadget("print")
	if g == nil {
		t.Fatalf("could not find [print] gadget")
	}
	_, ok := g.(*glow.Gadget)
	if !ok {
		t.Errorf("not a gadget: %v", g)
	}
}

func TestPrintGadget(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	b := &bytes.Buffer{}
	glow.Debug = b

	g := glow.NewGadget("print")
	g.Feed(0, glow.Message{"hello"})

	if b.String() != "hello\n" {
		t.Errorf("expected 'hello', got: %q", b)
	}
}

func TestPrintGadgetArg(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	b := &bytes.Buffer{}
	glow.Debug = b

	g := glow.NewGadget("print", 123)
	g.Feed(0, glow.Message{"hello"})

	if b.String() != "123 hello\n" {
		t.Errorf("expected '123 hello', got: %q", b)
	}
}

func TestPassGadgetExists(t *testing.T) {
	f, ok := glow.Registry["pass"]
	if !ok {
		t.Fatalf("could not find [pass] gadget")
	}
	g := f(nil)
	_, ok = g.(*glow.Gadget)
	if !ok {
		t.Errorf("not a gadget: %v", g)
	}
}

func TestPassAndPrintGadget(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	var b bytes.Buffer
	glow.Debug = &b

	g1 := glow.NewGadget("pass")
	g2 := glow.NewGadget("print")
	g1.Connect(0, g2, 0) // g1.out[0] => g2.in[0]

	g1.Feed(0, glow.Message{"howdy"})

	if b.String() != "howdy\n" {
		t.Errorf("expected 'howdy', got: %q", b)
	}
}

func TestBuildCircuit(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	b := &bytes.Buffer{}
	glow.Debug = b

	c := new(glow.Circuit)
	g := glow.NewGadget("pass")
	c.Add(g)
	c.Add(glow.NewGadget("print"))
	c.AddWire(0, 0, 1, 0)

	g.Feed(0, glow.Message{"bingo"})

	if b.String() != "bingo\n" {
		t.Errorf("expected 'bingo', got: %q", b)
	}
}

func TestCircuitInlet(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	b := &bytes.Buffer{}
	glow.Debug = b

	c := new(glow.Circuit)
	c.Add(glow.NewGadget("inlet"))
	c.Add(glow.NewGadget("print"))
	c.AddWire(0, 0, 1, 0)

	c.Feed(0, glow.Message{"foo"})

	if b.String() != "foo\n" {
		t.Errorf("expected 'foo', got: %q", b)
	}
}

func TestCircuitOutlet(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	b := &bytes.Buffer{}
	glow.Debug = b

	c := new(glow.Circuit)
	c.Add(glow.NewGadget("inlet"))
	c.Add(glow.NewGadget("outlet"))
	c.AddWire(0, 0, 1, 0)

	g := glow.NewGadget("print")
	c.Connect(0, g, 0) // c.out[0] => g.in[0]

	c.Feed(0, glow.Message{"bar"})

	if b.String() != "bar\n" {
		t.Errorf("expected 'bar', got: %q", b)
	}
}

func TestSwapGadget(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	b := &bytes.Buffer{}
	glow.Debug = b

	c := new(glow.Circuit)
	c.Add(glow.NewGadget("inlet"))
	c.Add(glow.NewGadget("swap", 123))
	c.Add(glow.NewGadget("print", 1))
	c.Add(glow.NewGadget("print", 2))
	c.AddWire(0, 0, 1, 0)
	c.AddWire(1, 0, 2, 0)
	c.AddWire(1, 1, 3, 0)

	c.Feed(0, glow.Message{111})
	c.Feed(0, glow.Message{222})

	if b.String() != "2 111\n1 123\n2 222\n1 123\n" {
		t.Errorf("expected 4 lines', got: %q", b)
	}
}

// this came straight out of the Pd-extended patch editor:
var swapPatch = `
#N canvas 673 402 450 300 10;
#X obj 75 101 swap 123;
#X obj 75 142 print 1;
#X obj 146 142 print 2;
#X obj 75 60 inlet;
#X connect 0 0 1 0;
#X connect 0 1 2 0;
#X connect 3 0 0 0;
`

func TestSwapPatch(t *testing.T) {
	tmp := glow.Debug
	defer func() { glow.Debug = tmp }()
	b := &bytes.Buffer{}
	glow.Debug = b

	c := glow.NewCircuitFromText(swapPatch)

	c.Feed(0, glow.Message{11})
	c.Feed(0, glow.Message{22})

	if b.String() != "2 11\n1 123\n2 22\n1 123\n" {
		t.Errorf("expected 4 lines', got: %q", b)
	}
}
