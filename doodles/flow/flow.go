// Flow implements a Pure-data like dataflow mechanism.
package flow

import (
	"fmt"
	"reflect"
)

var registry = map[string]func() Circuitry{}

// Circuitry is the collective name for gadgets and circuits.
type Circuitry interface {
	NumInlets() int
	Inlet(n int) *Inlet

	NumOutlets() int
	Outlet(n int) *Outlet

	Setup()
	Trigger()
	Cleanup()

	install(self Circuitry, name string, owner *Circuit) *Gadget
}

// Register a constructor for a named gadget type.
func Register(name string, f func() Circuitry) {
	registry[name] = f
}

// A circuit is a collection of gadgets.
type Circuit struct {
	Gadget
	gadgets map[string]*Gadget
}

// Add a new gadget to the circuit.
func (c *Circuit) Add(name, typ string) {
	g := registry[typ]()
	c.gadgets[name] = g.install(g, name, c)
}

// Add a new wire connection to a circuit.
func (c *Circuit) Connect(fname string, fpin int, tname string, tpin int) {
	fg := c.gadgets[fname]
	tg := c.gadgets[tname]
	fg.Outlet(fpin).Connect(tg.Inlet(tpin))
}

// Set a pin to a specified value.
func (c *Circuit) SetPin(name string, pin int, m Message) {
	g := c.gadgets[name]
	*g.Inlet(pin) = m
}

// Terminate all the gadgets in the circuit.
func (c *Circuit) Terminate() {
	for _, g := range c.gadgets {
		g.Terminate()
	}
	//close(c.feed)
	//<-c.done
}

// NewCircuit creates a new empty circuit.
func NewCircuit() *Circuit {
	return &Circuit{
		gadgets: make(map[string]*Gadget),
	}
}

// Trigger gets called when a message arrives at inlet zero.
func (c *Circuit) Trigger() {
}

// A Gadget is the building block for creating circuits with.
type Gadget struct {
	name    string
	owner   *Circuit
	feed    chan incoming
	done    chan struct{}
	inlets  []*Inlet
	outlets []*Outlet
}

// map inlets back to their owning gadgets for sending
// TODO will need a mutex unless this map becomes per-circuit
var inletMap = make(map[*Inlet]*Gadget)

// String returns the name of this gadget.
func (g *Gadget) String() string {
	return g.name
}

// Install intialises a gadget for use inside a circuit.
func (g *Gadget) install(self Circuitry, name string, owner *Circuit) *Gadget {
	g.name = name
	g.owner = owner
	g.feed = make(chan incoming)
	g.done = make(chan struct{})

	// use reflection to create lists of all the inlets and outlets
	gVal := reflect.ValueOf(self).Elem()
	gTyp := reflect.TypeOf(self).Elem()
	for i := 0; i < gVal.NumField(); i++ {
		fVal := gVal.Field(i)
		fTyp := gTyp.Field(i)
		fmt.Println("fp", i, fTyp.Name, fTyp.Type)
		switch fVal.Type().String() {
		case "flow.Inlet":
			in := fVal.Addr().Interface().(*Inlet)
			g.inlets = append(g.inlets, in)
			inletMap[in] = g
		case "flow.Outlet":
			out := fVal.Addr().Interface().(*Outlet)
			g.outlets = append(g.outlets, out)
		}
	}

	go g.run(self)

	return g
}

func (g *Gadget) run(self Circuitry) {
	defer func() {
		for _, x := range g.inlets {
			delete(inletMap, x)
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
	fmt.Println("Gadget setup:", g.name)
}

// Trigger gets called when a message arrives at inlet zero.
func (g *Gadget) Trigger() {
	fmt.Println("Gadget trigger:", g.name)
}

// Cleanup is called just after a gadget has finished normal processing.
func (g *Gadget) Cleanup() {
	fmt.Println("Gadget cleanup:", g.name)
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
	inletMap[i].feed <- incoming{m, i}
}

// An Outlet can be connected to zero or more inlets.
type Outlet []*Inlet

// FanOut returns the number of inlets currently connected.
func (o *Outlet) FanOut() int {
	return len(*o)
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
