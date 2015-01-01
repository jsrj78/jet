package main

import (
	"fmt"
	"reflect"
	"time"
)

type Circuitry interface {
	NumInlets() int
	NumOutlets() int
	Setup()
	Trigger()
	Cleanup()
}

type Circuit struct {
	Gadget
}

func NewCircuit() Circuitry {
	return &Circuit{}
}

func (c *Circuit) Trigger() {
}

type Gadget struct {
	inlets  []*Inlet
	outlets []*Outlet
	feed    chan Incoming
	done    chan interface{}
}

// map inlets back to their owning gadgets for sending
// TODO will need a mutex unless this map becomes per-circuit
var owners = make(map[*Inlet]*Gadget)

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

	g.feed = make(chan Incoming)
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

func (g *Gadget) Terminate() {
	close(g.feed)
	<-g.done
}

func (g *Gadget) NumInlets() int {
	return len(g.inlets)
}

func (g *Gadget) NumOutlets() int {
	return len(g.outlets)
}

func (g *Gadget) Setup() {
	fmt.Println("Gadget setup")
}

func (g *Gadget) Trigger() {
	fmt.Println("Gadget trigger")
}

func (g *Gadget) Cleanup() {
	fmt.Println("Gadget cleanup")
}

type Message interface{}

type Incoming struct {
	msg Message
	pin *Inlet
}

type Inlet Message

func SetInlet(i *Inlet, m Message) {
	owners[i].feed <- Incoming{m, i}
}

type Outlet []*Inlet

func (o *Outlet) IsActive() bool {
	return len(*o) > 0
}

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

func (o *Outlet) Connect(i *Inlet) {
	if o.indexOf(i) >= 0 {
		panic(fmt.Errorf("already connected"))
	}
	*o = append(*o, i)
}

func (o *Outlet) Disconnect(i *Inlet) {
	if n := o.indexOf(i); n >= 0 {
		*o = append((*o)[:n], (*o)[n+1:]...)
	}
}

type MetroG struct {
	Gadget
	Out Outlet
}

func (g *MetroG) Setup() {
	fmt.Println("MetroG setup")
	time.Sleep(time.Second)
	g.Out.Send("hi!")
	time.Sleep(time.Second)
	g.Out.Send("ha!")
	time.Sleep(time.Second)
	g.Out.Send("ho!")
}

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
	g2 := new(PrintG)
	g2.Launch(g2)
	fmt.Println("g1#out", g1.NumOutlets(), "g2#in", g2.NumInlets())

	g1.Out.Connect(&g2.In)

	g1.Terminate()
	g2.Terminate()
	fmt.Println("exit", len(owners))
}
