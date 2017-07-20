// The Glow package implements a dataflow engine in Go.
// It was inspired by Pure Data (http://puredata.info).
package glow

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// NewMsg constructs a new message object.
func NewMsg(args ...interface{}) Msg {
	return Msg(args)
}

// A Msg is what gets passed around: a "bang", int, string, or vector.
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
		_, e := strconv.Atoi(s)
		if len(s) == 0 {
			s = `""`
		} else if e == nil || len(t) != len(s)+2 || strings.Contains(s, " ") {
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
		mi := m[index]
		if mi == nil {
			m = Msg{}
		} else if m2, ok := mi.(Msg); ok {
			m = m2
		} else {
			m = NewMsg(mi)
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

// The Registry is a collection of named gadget constructors.
var Registry = map[string]func(args Msg) Gadgetry{}

// Gadgetry is the common interface for all gadgets and circuits.
type Gadgetry interface {
	AddedTo(*Circuit)
	Connect(int, Gadgetry, int)
	Feed(int, Msg)
	Emit(int, Msg)
}

// A Gadget is the base type for all gadgets.
type Gadget struct {
	inlets  []Inlet
	outlets []Outlet
	onAdded func(*Circuit)
}

// An Endpoint is a reference to a specific inlet or outlet in a gadget.
type Endpoint struct {
	gadget Gadgetry
	index  int
}

// An Inlet is an endpoint which accepts messages.
type Inlet struct {
	handler func(m Msg)
}

// An Outlet is an endpoint which publishes messages.
type Outlet []Endpoint

// NewGadget instantiates a gadget from the registry, with optional args.
func NewGadget(args ...interface{}) Gadgetry {
	name := args[0].(string)
	r, ok := Registry[name]
	if !ok {
		//fmt.Println("unknown gadget:", args)
		return nil
	}
	return r(Msg(args[1:]))
}

// AddInlet sets up a new gadget inlet.
func (g *Gadget) AddInlet(f func(m Msg)) {
	g.inlets = append(g.inlets, Inlet{handler: f})
}

// AddOutlets sets up new gadget outlets.
func (g *Gadget) AddOutlets(n int) int {
	i := len(g.outlets)
	g.outlets = append(g.outlets, make([]Outlet, n)...)
	return i
}

// AddedTo is called when a gadget has been added to a circuit.
func (g *Gadget) AddedTo(c *Circuit) {
	if g.onAdded != nil {
		g.onAdded(c)
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

// A Circuit is a composition of gadgets, including sub-circuits.
type Circuit struct {
	Gadget
	gadgets []Gadgetry
}

// Add a new gadget (or sub-circuit) to a circuit.
func (c *Circuit) Add(g Gadgetry) {
	c.gadgets = append(c.gadgets, g)
	g.AddedTo(c)
}

// AddWire adds a connection from one gadget's outlet to another's inlet.
func (c *Circuit) AddWire(srcg, srco, dstg, dsti int) {
	c.gadgets[srcg].Connect(srco, c.gadgets[dstg], dsti)
}

// MsgFromString parses a string and returns a message constructed from it.
func MsgFromString(s string) Msg {
	m := Msg{}
	for _, x := range strings.Split(s, " ") {
		if v, e := strconv.Atoi(x); e == nil {
			m = append(m, v)
		} else {
			m = append(m, x)
		}
	}
	return m
}

// NewCircuitFromText constructs a circuit from a Pd text representation.
func NewCircuitFromText(text string) Gadgetry {
	c := new(Circuit)
	for _, s := range strings.Split(text, "\n") {
		if strings.HasPrefix(s, "#X ") && strings.HasSuffix(s, ";") {
			m := MsgFromString(s[3 : len(s)-1])
			switch m[0] {
			case "obj":
				c.Add(NewGadget(m[3:]...))
			case "connect":
				c.AddWire(m[1].(int), m[2].(int), m[3].(int), m[4].(int))
			}
		}
	}
	return c
}

// An Event represents the state of a single registered event handler.
type Event struct {
	callback func(Msg)
	topic    string
	period   time.Duration
}

// An EventEmitter can be used as local pubsub mechanism.
type EventEmitter map[string][]*Event

// On subscribes to a specific topic.
func (ee EventEmitter) On(s string, f func(Msg)) *Event {
	e := &Event{callback: f, topic: s, period: 0}
	handlers, _ := ee[s]
	ee[s] = append(handlers, e)
	return e
}

// Emit triggers the specified topic.
func (ee EventEmitter) Emit(s string, args ...interface{}) {
	handlers, _ := ee[s]
	for _, e := range handlers {
		e.callback(NewMsg(args...))
	}
}
