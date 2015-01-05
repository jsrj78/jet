package flow

import (
	"fmt"
	"reflect"
)

// Gadget is the building block for creating circuits.
type Gadget struct {
	owner   *Circuit
	self    Circuitry
	feed    chan incoming
	done    chan struct{}
	inlets  []*Inlet
	outlets []*Outlet
}

// install intialises a gadget for use inside a circuit.
func (g *Gadget) install(self Circuitry, owner *Circuit) *Gadget {
	g.owner = owner
	g.self = self
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

	go g.run()

	return g
}

func (g *Gadget) run() {
	defer g.unlink()

	g.self.Setup()
	for x := range g.feed {
		if x.pin == nil {
			g.self.Control(x.msg)
		} else {
			*x.pin = x.msg
			if x.pin == g.inlets[0] {
				g.self.Trigger()
			}
		}
	}
	g.self.Cleanup()
}

// unlink from inletMap and from all outlets connected to this gadget.
func (g *Gadget) unlink() {
	for _, x := range g.inlets {
		delete(inletMap, x)
		// TODO inefficient, this iterates over all possible combinations
		// could set up a map with all *relevant* inlets or outlets instead
		for _, y := range g.owner.gadgets {
			for i := 0; i < y.NumOutlets(); i++ {
				y.Outlet(i).Disconnect(x)
			}
		}
	}
	close(g.done)
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

// Control gets called with messages sent to the special nil inlet.
func (g *Gadget) Control(m Message) {
	fmt.Println("Gadget control:", m)
}

// Trigger gets called when a message arrives at inlet zero.
func (g *Gadget) Trigger() {
	fmt.Println("Gadget trigger")
}

// Cleanup is called just after a gadget has finished normal processing.
func (g *Gadget) Cleanup() {
	fmt.Println("Gadget cleanup")
}
