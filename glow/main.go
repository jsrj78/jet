package main

import (
	"fmt"
	"io"
	"os"
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
	In(i int, m Msg)
}

// Gadget is the base type for all gadgets.
type Gadget struct {
	inlets []func(m Msg)
}

// AddInlet is used to set up all inlets.
func (g *Gadget) AddInlet(f func(m Msg)) {
	g.inlets = append(g.inlets, f)
}

// In sends a message to a specific inlet (indexed from 0 upwards).
func (g *Gadget) In(i int, m Msg) {
	g.inlets[i](m)
}

func init() {
	Registry["print"] = func() Gadgetry {
		g := new(Gadget)
		g.AddInlet(func(m Msg) {
			Debug.Write([]byte(m.AsString()))
		})
		return g
	}
}
