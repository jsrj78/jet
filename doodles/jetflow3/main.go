// Main "jetflow3" loads all the generic code from a separate "flow" package.
package main

import (
	"../flow"
	"fmt"
	"time"
)

// sample gadgets for a trivial pipeline: MetroG -> RepeatG -> PrintG

func init() {
	flow.Register("metro", func() flow.Circuitry { return new(MetroG) })
	flow.Register("repeat", func() flow.Circuitry { return new(RepeatG) })
	flow.Register("print", func() flow.Circuitry { return new(PrintG) })
}

// A MetroG gadget sends out periodic messages.
type MetroG struct {
	flow.Gadget
	Out flow.Outlet
}

func (g *MetroG) Setup() {
	// TODO this is test code, needs a real implementation
	fmt.Println("MetroG setup")
	g.Gadget.Setup()
	time.Sleep(500 * time.Millisecond)
	g.Out.Send("hi!")
	time.Sleep(500 * time.Millisecond)
	g.Out.Send("ha!")
	time.Sleep(500 * time.Millisecond)
	g.Out.Send("ho!")
}

// A RepeatG gadget repeats each incoming message Num times.
type RepeatG struct {
	flow.Gadget
	In  flow.Inlet
	Num flow.Inlet
	Out flow.Outlet
}

func (g *RepeatG) Trigger() {
	for i := 0; i < g.Num.(int); i++ {
		g.Out.Send(g.In)
	}
}

// A PrintG gadget prints everything received on its main inlet.
type PrintG struct {
	flow.Gadget
	In flow.Inlet
}

func (g *PrintG) Trigger() {
	fmt.Println(g.In)
}

func main() {
	fmt.Println("jetflow 0.2.4")

	c := flow.NewCircuit()
	c.Add("g1", "metro")
	c.Add("g2", "repeat")
	c.Add("g3", "print")

	c.SetPin("g2", 1, 3)

	c.Connect("g1", 0, "g2", 0)
	c.Connect("g2", 0, "g3", 0)

	c.Terminate()
	fmt.Println("exit")
}
