// Flow implements a Pure-data like dataflow mechanism.
package flow

import (
	"fmt"
)

var registry = map[string]func() Circuitry{}

// Circuitry is the collective name for gadgets and circuits.
type Circuitry interface {
	NumInlets() int
	Inlet(n int) *Inlet

	NumOutlets() int
	Outlet(n int) *Outlet

	Setup()
	Loop()
	Trigger()
	Cleanup()

	install(self Circuitry, name string, owner *Circuit) *Gadget
}

// Register a constructor for a named gadget type.
func Register(name string, f func() Circuitry) {
	registry[name] = f
}

// map inlets back to their owning gadgets for sending
// TODO will need a mutex or channel, see sendToInlet()
var inletMap = make(map[*Inlet]*Gadget)

// A message is a generic data item which can be sent between gadgets.
type Message interface{}

// An Inlet is a slot to store incoming messages.
type Inlet Message

type incoming struct {
	pin *Inlet
	msg Message
}

// SendToInlet will store a message into a specified inlet.
func sendToInlet(i *Inlet, m Message) {
	// TODO access to inletMap is not protected against map changes right now
	// could be either a mutex or an additional global channel added in front
	inletMap[i].feed <- incoming{pin: i, msg: m}
}

// An Outlet can be connected to zero or more inlets.
type Outlet []*Inlet

// FanOut returns the number of inlets currently connected.
func (o *Outlet) FanOut() int {
	return len(*o)
}

// Send will send out a message to all the attached inlets.
func (o *Outlet) Send(m Message) {
	// TODO add logging capability
	// could use an "outletMap" to retrieve the sending gadget's name
	for _, x := range *o {
		sendToInlet(x, m)
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
