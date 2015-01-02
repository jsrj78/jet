package main

import (
	"fmt"
	"reflect"
	"time"
)

// Circuitry is the collective name for gadgets and circuits.
type Circuitry interface {
	NumInlets() int
	Inlet(n int) *Inlet

	NumOutlets() int
	Outlet(n int) *Outlet

	Setup()
	Trigger()
	Cleanup()
}

// A circuit is a collection of gadgets.
type Circuit struct {
	Gadget
	children map[string]*Gadget
}

// Add a gadget to the circuit.
func (c *Circuit) Add(name string, g *Gadget) {
	c.children[name] = g
}

// NewCircuit creates a new empty circuit.
func NewCircuit() Circuitry {
	return &Circuit{}
}

// Trigger gets called when a message arrives at inlet zero.
func (c *Circuit) Trigger() {
}

// A Gadget is the building block for creating circuits with.
type Gadget struct {
	inlets  []*Inlet
	outlets []*Outlet
	feed    chan incoming
	done    chan interface{}
}

// map inlets back to their owning gadgets for sending
// TODO will need a mutex unless this map becomes per-circuit
var owners = make(map[*Inlet]*Gadget)

// Launch must be called once to intialise a gadget.
func (g *Gadget) Launch(self Circuitry) {
	// use reflection to create lists of all the inlets and outlets
	gVal := reflect.ValueOf(self).Elem()
	gTyp := reflect.TypeOf(self).Elem()
	for i := 0; i < gVal.NumField(); i++ {
		fVal := gVal.Field(i)
		fTyp := gTyp.Field(i)
		fmt.Println("fp", i, fTyp.Name, fTyp.Type)
		switch fVal.Type().String() {
		case "main.Inlet":
			in := fVal.Addr().Interface().(*Inlet)
			g.inlets = append(g.inlets, in)
			owners[in] = g
		case "main.Outlet":
			out := fVal.Addr().Interface().(*Outlet)
			g.outlets = append(g.outlets, out)
		}
	}

	g.feed = make(chan incoming)
	g.done = make(chan interface{})

	go g.runner(self)
}

func (g *Gadget) runner(self Circuitry) {
	defer func() {
		for _, x := range g.inlets {
			delete(owners, x)
		}
		close(g.done)
	}()

	self.Setup()
	for x := range g.feed {
		*x.pin = x.msg
		if x.pin == g.inlets[0] {
			self.Trigger()
		}
	}
	self.Cleanup()
}

// Terminate causes the gadget to end and cleanup, and returns when it's done.
func (g *Gadget) Terminate() {
	close(g.feed)
	<-g.done
}

// NumInlets returns the number of inlets in this gadget.
func (g *Gadget) NumInlets() int {
	return len(g.inlets)
}

// Inlet returns a pointer to the n'th inlet in this gadget.
func (g *Gadget) Inlet(n int) *Inlet {
	return g.inlets[n]
}

// NumOutlets returns the number of outlets in this gadget.
func (g *Gadget) NumOutlets() int {
	return len(g.outlets)
}

// Outlet returns a pointer to the n'th outlet in this gadget.
func (g *Gadget) Outlet(n int) *Outlet {
	return g.outlets[n]
}

// Setup is called just before a gadget starts normal processing.
func (g *Gadget) Setup() {
	fmt.Println("Gadget setup")
}

// Trigger gets called when a message arrives at inlet zero.
func (g *Gadget) Trigger() {
	fmt.Println("Gadget trigger")
}

// Cleanup is called just after a gadget has finished normal processing.
func (g *Gadget) Cleanup() {
	fmt.Println("Gadget cleanup")
}

// A message is a generic data item which can be sent between gadgets.
type Message interface{}

type incoming struct {
	msg Message
	pin *Inlet
}

// An Inlet is a slot to store incoming messages.
type Inlet Message

// SetInlet will store a message into a specified inlet.
func SetInlet(i *Inlet, m Message) {
	owners[i].feed <- incoming{m, i}
}

// An Outlet can be connected to zero or more inlets.
type Outlet []*Inlet

// IsActive returns true if the outlet is connected to at least one inlet.
func (o *Outlet) IsActive() bool {
	return len(*o) > 0
}

// Send will send out a message to all the attached inlets.
func (o *Outlet) Send(m Message) {
	for _, x := range *o {
		SetInlet(x, m)
	}
}

func (o *Outlet) indexOf(i *Inlet) int {
	for n, x := range *o {
		if x == i {
			return n
		}
	}
	return -1
}

// Connect an outlet to a specified inlet.
func (o *Outlet) Connect(i *Inlet) {
	if o.indexOf(i) >= 0 {
		panic(fmt.Errorf("already connected"))
	}
	*o = append(*o, i)
}

// Disconnect a specified inlet from the outlet.
func (o *Outlet) Disconnect(i *Inlet) {
	if n := o.indexOf(i); n >= 0 {
		*o = append((*o)[:n], (*o)[n+1:]...)
	}
}

// sample gadgets for a trivial pipeline: MetroG -> RepeatG -> PrintG

// A MetroG gadget sends out periodic messages.
type MetroG struct {
	Gadget
	Out Outlet
}

func (g *MetroG) Setup() {
	// TODO this is test code, needs a real implementation
	fmt.Println("MetroG setup")
	time.Sleep(time.Second)
	g.Out.Send("hi!")
	time.Sleep(time.Second)
	g.Out.Send("ha!")
	time.Sleep(time.Second)
	g.Out.Send("ho!")
}

// A RepeatG gadget repeats each incoming message Num times.
type RepeatG struct {
	Gadget
	In  Inlet
	Num Inlet
	Out Outlet
}

func (g *RepeatG) Trigger() {
	for i := 0; i < g.Num.(int); i++ {
		g.Out.Send(g.In)
	}
}

// A PrintG gadget prints everything received on its main inlet.
type PrintG struct {
	Gadget
	In Inlet
}

func (g *PrintG) Trigger() {
	fmt.Println(g.In)
}

func main() {
	fmt.Println("jetflow 0.2.2")
	g1 := new(MetroG)
	g1.Launch(g1)
	g2 := new(RepeatG)
	g2.Launch(g2)
	g3 := new(PrintG)
	g3.Launch(g3)

	g2.Num = 3

	g1.Out.Connect(&g2.In)
	g2.Out.Connect(&g3.In)

	g1.Terminate()
	g2.Terminate()
	g3.Terminate()
	fmt.Println("exit", len(owners))
}
