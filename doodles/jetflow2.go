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

var inlets = make(map[*Inlet]*Gadget)

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
			iPtr := fVal.Addr().Interface().(*Inlet)
			g.inlets = append(g.inlets, iPtr)
			inlets[iPtr] = g
		case "main.Outlet":
			g.outlets = append(g.outlets, fVal.Addr().Interface().(*Outlet))
		}
	}

	g.feed = make(chan Incoming)
	g.done = make(chan interface{})

	go g.runner(self)
}

func (g *Gadget) runner(self Circuitry) {
	self.Setup()
	for x := range g.feed {
		*x.pin = Inlet{x.msg}
		if x.pin == g.inlets[0] {
			self.Trigger()
		}
	}
	self.Cleanup()
	close(g.done)
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

type Inlet struct {
	Message
}

func (i *Inlet) Set(m Message) {
	inlets[i].feed <- Incoming{m, i}
}

type Outlet []*Inlet

func (o *Outlet) IsActive() bool {
	return len(*o) > 0
}

func (o *Outlet) Send(m Message) {
	for _, x := range *o {
		x.Set(m)
	}
}

func Connect(from Circuitry, fpin int, to Circuitry, tpin int) {
	//*f.outlets[fpin] = append(*f.outlets[fpin], to.inlets[tpin])
}

func Disconnect(from *Gadget, fpin int, to *Gadget, tpin int) {
	// ...
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
	fmt.Println("jethub 0.2.1")
	g1 := new(MetroG)
	g1.Launch(g1)
	g2 := new(PrintG)
	g2.Launch(g2)
	fmt.Println("g1#out", g1.NumOutlets(), "g2#in", g2.NumInlets())

	//Connect(g1, 0, g2, 0)
	//*f.outlets[fpin] = append(*f.outlets[fpin], to.inlets[tpin])
	// FIXME couldn't get Connect to work yet, so wiring is hard-coded for now
	*g1.outlets[0] = append(*g1.outlets[0], g2.inlets[0])

	g1.Terminate()
	g2.Terminate()
	fmt.Println("exit")
}
