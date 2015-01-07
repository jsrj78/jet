package main

import (
	"fmt"
)

type Circuitry interface {
	NumInlets() int
	NumOutlets() int
}

type Gadget struct {
	inlets  []*Inlet
	outlets []*Outlet
}

func (g *Gadget) NumInlets() int {
	return len(g.inlets)
}

func (g *Gadget) NumOutlets() int {
	return len(g.outlets)
}

func (g *Gadget) Trigger() {
	fmt.Println("triggered!")
}

type Message interface{}

type Inlet struct {
	value Message
	owner *Gadget
}

func (i *Inlet) Set(m Message) {
	i.value = m
	if i == i.owner.inlets[0] {
		i.owner.Trigger()
	}
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

func Connect(from *Gadget, fpin int, to *Gadget, tpin int) {
	// ...
}

func Disconnect(from *Gadget, fpin int, to *Gadget, tpin int) {
	// ...
}

func main() {
	fmt.Println("jethub 0.1")
}
