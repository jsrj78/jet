// Main "jetflow3" loads all the generic code from a separate "flow" package.
package main

import (
	"fmt"
	"time"

	"github.com/jeelabs/jet/flow"
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

	// TODO generate these requests from text lines with a suitable parser
	c.Request("add", "metro")
	c.Request("add", "repeat")
	c.Request("add", "print")
	c.Request("connect", 0, 0, 1, 0)
	c.Request("connect", 1, 0, 2, 0)
	c.Request("send2pin", 1, 1, 3)

	c.Terminate()
	fmt.Println("exit")
}
