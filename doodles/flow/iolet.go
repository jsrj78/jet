package flow

import (
	"fmt"
)

// map inlets back to their owning gadgets for sending
// TODO will need a mutex or channel, see sendToInlet()
var inletMap = map[*Inlet]*Gadget{}

// incoming is used to store a message which needs to be sent to an inlet.
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

// An Inlet is a slot to store incoming messages.
type Inlet Message

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

// indexOf returns the index of an inlet in the outlet's list, or -1.
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
