package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fmt.Println("hello")
}

// NewMsg constructs a new message object.
func NewMsg(args ...interface{}) Msg {
	return Msg(args)
}

// Msg is what gets passed around: a "bang", int, string, or vector.
type Msg []interface{}

// String returns a nice string representation of a message.
func (m Msg) String() string {
	if m.IsBang() {
		return "[]"
	}
	if m.IsInt() {
		return fmt.Sprint(m.AsInt())
	}
	if m.IsString() {
		s := m.AsString()
		t := fmt.Sprintf("%q", s)
		if len(s) == 0 {
			s = `""`
		} else if len(t) != len(s)+2 || strings.Contains(s, " ") {
			s = t
		}
		return s
	}
	v := []string{}
	for i := range m {
		e := m.At(i)
		s := e.String()
		if !e.IsBang() && !e.IsInt() && !e.IsString() {
			s = "[" + s + "]"
		}
		v = append(v, s)
	}
	return strings.Join(v, " ")
}

// At indexes arbitrarily-deeply-nested message structures.
func (m Msg) At(indices ...int) Msg {
	for _, index := range indices {
		if index >= len(m) {
			return Msg{}
		}
		if m2, ok := m[index].(Msg); ok {
			m = m2
		} else {
			m = NewMsg(m[index])
		}
	}
	return m
}

// IsBang returns true if m is a "bang".
func (m Msg) IsBang() bool {
	return len(m) == 0
}

// IsInt returns true if m is an int.
func (m Msg) IsInt() (ok bool) {
	if len(m) == 1 {
		_, ok = m[0].(int)
	}
	return
}

// IsString returns true if m is a string.
func (m Msg) IsString() (ok bool) {
	if len(m) == 1 {
		_, ok = m[0].(string)
	}
	return
}

// AsInt returns the int in m, else 0.
func (m Msg) AsInt() int {
	if m.IsInt() {
		return m[0].(int)
	}
	//fmt.Println("not an int:", m)
	return 0
}

// AsString returns the string in m, else "".
func (m Msg) AsString() string {
	if m.IsString() {
		return m[0].(string)
	}
	//fmt.Println("not a string:", m)
	return ""
}

// Debug is a Writer for debugging output.
var Debug io.Writer = os.Stdout

// Registry is a named collection of gadgets.
var Registry = map[string]func() Gadgetry{}

// Gadgetry is the common interface to gadgets.
type Gadgetry interface {
	Connect(o int, d Gadgetry, i int)
	Feed(i int, m Msg)
	Emit(i int, m Msg)
}

// Gadget is the base type for all gadgets.
type Gadget struct {
	inlets  []Inlet
	outlets []Outlet
}

// Endpoint is a reference to a specific inlet or outlet in a gadget.
type Endpoint struct {
	gadget Gadgetry
	index  int
}

// Inlet is an endpoint which accepts messages.
type Inlet struct {
	handler func(m Msg)
}

// Outlet is an endpoint which publishes messages.
type Outlet []Endpoint

// AddInlet is used to set up each inlet.
func (g *Gadget) AddInlet(f func(m Msg)) {
	g.inlets = append(g.inlets, Inlet{handler: f})
}

// AddOutlets is used to set up new outlets.
func (g *Gadget) AddOutlets(n int) {
	for i := 0; i < n; i++ {
		g.outlets = append(g.outlets, Outlet{})
	}
}

// Connect adds a connection from a gadget output to a gadget input.
func (g *Gadget) Connect(o int, d Gadgetry, i int) {
	g.outlets[o] = append(g.outlets[o], Endpoint{d, i})
}

// Feed accepts a message for a specific inlet (indexed from 0 upwards).
func (g *Gadget) Feed(i int, m Msg) {
	g.inlets[i].handler(m)
}

// Emit sends a message to a specific outlet (indexed from 0 upwards).
func (g *Gadget) Emit(o int, m Msg) {
	for _, ep := range g.outlets[o] {
		ep.gadget.Feed(ep.index, m)
	}
}
