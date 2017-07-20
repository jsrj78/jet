// The Glow package implements a dataflow engine in Go.
// It was inspired by Pure Data (http://puredata.info).
package glow

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// A Message is what gets passed around: a "bang", int, string, or vector.
type Message []interface{}

// String returns a nice string representation of a message.
func (m Message) String() string {
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
func (m Message) At(indices ...int) Message {
	for _, index := range indices {
		if index >= len(m) {
			return nil
		}
		mi := m[index]
		if mi == nil {
			return nil
		}
		if m2, ok := mi.(Message); ok {
			m = m2
		} else {
			m = Message{mi}
		}
	}
	return m
}

// IsBang returns true if m is a "bang".
func (m Message) IsBang() bool {
	return len(m) == 0
}

// IsInt returns true if m is an int.
func (m Message) IsInt() (ok bool) {
	if len(m) == 1 {
		_, ok = m[0].(int)
	}
	return
}

// IsString returns true if m is a string.
func (m Message) IsString() (ok bool) {
	if len(m) == 1 {
		_, ok = m[0].(string)
	}
	return
}

// AsInt returns the int in m, else 0.
func (m Message) AsInt() int {
	if m.IsInt() {
		return m[0].(int)
	}
	//fmt.Println("not an int:", m)
	return 0
}

// AsString returns the string in m, else "".
func (m Message) AsString() string {
	if m.IsString() {
		return m[0].(string)
	}
	//fmt.Println("not a string:", m)
	return ""
}

// Debug is a Writer for debugging output.
var Debug io.Writer = os.Stdout

// The Registry is a collection of named gadget constructors.
var Registry = map[string]func(args Message) Gadgetry{}

// Gadgetry is the common interface for all gadgets and circuits.
type Gadgetry interface {
	AddedTo(*Circuit)
	Connect(int, Gadgetry, int)
	Feed(int, Message)
	Emit(int, Message)
}

// A Gadget is the base type for all gadgets.
type Gadget struct {
	OnAdded func(*Circuit) // called when we've been added to a circuit

	ins  []Inlet
	outs []Outlet
}

// An endpoint is a reference to a specific inlet or outlet in a gadget.
type endpoint struct {
	gadget Gadgetry
	index  int
}

// An Inlet is an endpoint which accepts messages.
type Inlet struct {
	handler func(m Message)
}

// An Outlet is an endpoint which publishes messages.
type Outlet []endpoint

// NewGadget creates a new gadget with default settings.
func NewGadget() *Gadget {
	return new(Gadget)
}

// LookupGadget instantiates a gadget from the registry, with optional args.
func LookupGadget(name string, args ...interface{}) Gadgetry {
	r, ok := Registry[name]
	if !ok {
		//fmt.Println("unknown gadget:", args)
		return nil
	}
	return r(args)
}

// AddInlet sets up a new gadget inlet.
func (g *Gadget) AddInlet(f func(m Message)) {
	g.ins = append(g.ins, Inlet{handler: f})
}

// AddOutlets sets up new gadget outlets.
func (g *Gadget) AddOutlets(n int) int {
	i := len(g.outs)
	g.outs = append(g.outs, make([]Outlet, n)...)
	return i
}

// AddedTo is called when a gadget has been added to a circuit.
func (g *Gadget) AddedTo(c *Circuit) {
	if g.OnAdded != nil {
		g.OnAdded(c)
	}
}

// Connect adds a connection from a gadget output to a gadget input.
func (g *Gadget) Connect(o int, d Gadgetry, i int) {
	g.outs[o] = append(g.outs[o], endpoint{d, i})
}

// Feed accepts a message for a specific inlet (indexed from 0 upwards).
func (g *Gadget) Feed(i int, m Message) {
	g.ins[i].handler(m)
}

// Emit sends a message to a specific outlet (indexed from 0 upwards).
func (g *Gadget) Emit(o int, m Message) {
	for _, ep := range g.outs[o] {
		ep.gadget.Feed(ep.index, m)
	}
}

// A Circuit is a composition of gadgets, including sub-circuits.
type Circuit struct {
	Gadget
	Notifier

	gadgets []Gadgetry
}

// NewCircuit creates a new empty circuit
func NewCircuit() *Circuit {
	c := new(Circuit)
	c.Notifier = make(Notifier)
	return c
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

// ParseAsMessage parses a string and returns a message constructed from it.
func ParseAsMessage(s string) (m Message) {
	for _, x := range strings.Split(s, " ") {
		if v, e := strconv.Atoi(x); e == nil {
			m = append(m, v)
		} else {
			m = append(m, x)
		}
	}
	return
}

// NewCircuitFromText constructs a circuit from a Pd text representation.
func NewCircuitFromText(text string) Gadgetry {
	c := NewCircuit()
	for _, s := range strings.Split(text, "\n") {
		if strings.HasPrefix(s, "#X ") && strings.HasSuffix(s, ";") {
			m := ParseAsMessage(s[3 : len(s)-1])
			switch m[0] {
			case "obj":
				c.Add(LookupGadget(m[3].(string), m[4:]...))
			case "connect":
				c.AddWire(m[1].(int), m[2].(int), m[3].(int), m[4].(int))
			}
		}
	}
	return c
}

// A listener responds to notifications.
type listener struct {
	callback func(Message)
	topic    string
	periodic bool
}

// A Notifier calls listeners interested in a topic or after a timeout.
type Notifier map[string][]*listener

// On subscribes to a specific topic.
func (nf Notifier) On(s string, f func(Message)) *listener {
	e := &listener{callback: f, topic: s}
	lv, _ := nf[s]
	nf[s] = append(lv, e)
	return e
}

// Off unsubscribes an existing listener.
func (nf Notifier) Off(l *listener) {
	lv := nf[l.topic]
	for i, x := range lv {
		if l == x {
			n := copy(lv[i+1:], lv[i:])
			lv = lv[:i+n]
		}
	}
	if len(lv) > 0 {
		nf[l.topic] = lv
	} else {
		delete(nf, l.topic)
	}
}

// Notify informs all listeners of a specific topic.
func (nf Notifier) Notify(s string, args ...interface{}) {
	l, _ := nf[s]
	for _, e := range l {
		e.callback(args)
	}
}

// Now is the current time, either real or simulated.
var Now int

// The timers map keeps track of all global timeout listeners.
var timers = make(Notifier)

// NextTimer is set to the lowest pending timer value.
var NextTimer = -1

// TODO could use a string and fixed-length numbers to avoid many conversions
// listeners could even be combined with another notifier if prefixed by "t:"

// Run advances (real or simulated) time and triggers all timers as scheduled.
func Run(ms int) {
	tlimit := Now + ms
	for NextTimer >= 0 && NextTimer <= tlimit {
		Now = NextTimer                // this is where simulated time advances
		timers.Notify(fmt.Sprint(Now)) // fire all the matching pending timers
		lookForNextTimer()             // figure out when next timer must run
	}
	Now = tlimit // final time jump
}

// lookForNextTimer scans the timers to find the first pending one.
func lookForNextTimer() {
	NextTimer = -1
	for topic := range timers {
		t, _ := strconv.Atoi(topic)
		fixNextTimer(t)
	}
}

// adjustNextTimer updates the next timeout we need to process.
func fixNextTimer(t int) {
	if t < NextTimer || NextTimer < 0 {
		NextTimer = t
	}
}

// use a separate type for timer listeners to avoid mixups
type timer listener

// SetTimer schedules a one-shot notification.
func SetTimer(ms int, f func()) *timer {
	tsched := Now + ms
	var t *timer
	t = (*timer)(timers.On(fmt.Sprint(tsched), func(Message) {
		CancelTimer(t)
		if t.periodic {
			SetPeriodic(ms, f)
		}
		f()
	}))
	fixNextTimer(tsched)
	return t
}

// SetPeriodic schedules a repeating notification.
func SetPeriodic(ms int, f func()) *timer {
	t := SetTimer(ms, f)
	t.periodic = true
	return t
}

// CancelTimer drops a pending timer notification.
func CancelTimer(t *timer) {
	timers.Off((*listener)(t))
	// make sure NextTimer remains valid
	t1, _ := strconv.Atoi(t.topic)
	if t1 == NextTimer {
		lookForNextTimer()
	}
}
